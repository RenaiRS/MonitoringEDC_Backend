package controllers

import (
	"errors"

	"gorm.io/gorm"

	"backend/database"
	"backend/models"
	"backend/models/dto"
	"backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

func UploadVendorExcel(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "File tidak ditemukan"})
	}

	if !utils.IsExcelFile(file.Filename) {
		return c.Status(400).JSON(fiber.Map{"message": "File harus berformat Excel"})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Gagal membuka file"})
	}
	defer src.Close()

	xl, err := excelize.OpenReader(src)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Gagal membaca file excel"})
	}

	sheet := xl.GetSheetName(0)
	rows, err := xl.GetRows(sheet)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Gagal membaca sheet"})
	}

	var inserted, updated, skipped int
	var skipReasons = make(map[string]int)

	for i := 1; i < len(rows); i++ {
		if utils.IsVendorHeaderRow(rows[i]) {
			skipped++
			skipReasons["header"]++
			continue
		}

		if len(rows[i]) == 0 || utils.IsEmptyRow(rows[i]) {
			skipped++
			skipReasons["empty_row"]++
			continue
		}

		tid := dto.GetRawCell(xl, sheet, i, 1)
		mid := dto.GetRawCell(xl, sheet, i, 2)
		kota := dto.GetRawCell(xl, sheet, i, 6)
		cabang := dto.GetRawCell(xl, sheet, i, 7)
		tipeEDC := dto.GetRawCell(xl, sheet, i, 8)

		if tid == "" {
			skipped++
			skipReasons["no_tid"]++
			continue
		}

		var mesin models.MesinEDC
		err := database.DB.Where("terminal_id = ?", tid).First(&mesin).Error

		if err == nil {
			vendorData := models.MesinEDC{
				MID:     mid,
				Kota:    kota,
				Cabang:  cabang,
				TipeEDC: tipeEDC,
			}

			updates := utils.MergeVendorData(&mesin, vendorData)

			if len(updates) == 0 {
				skipped++
				skipReasons["no_changes"]++
				continue
			}

			if err := database.DB.Model(&mesin).Updates(updates).Error; err != nil {
				skipped++
				skipReasons["update_failed"]++
				continue
			}
			var sewa models.Sewa
			errSewa := database.DB.Where("mesin_id = ?", mesin.ID).First(&sewa).Error

			if errors.Is(errSewa, gorm.ErrRecordNotFound) {
				sewa = models.Sewa{
					MesinID:      mesin.ID,
					StatusSewa:   "berakhir",
					BiayaBulanan: 0,
				}
				database.DB.Create(&sewa)
			}

			updated++
			continue
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newMesin := models.MesinEDC{
				TerminalID:  tid,
				MID:         mid,
				Kota:        kota,
				Cabang:      cabang,
				TipeEDC:     tipeEDC,
				StatusData:  "vendor_only",
				StatusMesin: "aktif",
				LetakMesin:  "vendor",
			}

			if err := database.DB.Create(&newMesin).Error; err != nil {
				skipped++
				skipReasons["create_failed"]++
				continue
			}

			sewa := models.Sewa{
				MesinID:      newMesin.ID,
				StatusSewa:   "berakhir",
				BiayaBulanan: 0,
			}

			database.DB.Create(&sewa)

			inserted++
			continue
		}

		skipped++
		skipReasons["unknown_error"]++
	}

	return c.JSON(fiber.Map{
		"message":      "Upload vendor selesai",
		"inserted":     inserted,
		"updated":      updated,
		"skipped":      skipped,
		"skip_reasons": skipReasons,
		"total_rows":   len(rows) - 1,
	})
}

func UploadBankExcel(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "File tidak ditemukan"})
	}

	if !utils.IsExcelFile(file.Filename) {
		return c.Status(400).JSON(fiber.Map{"message": "File harus berformat Excel"})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Gagal membuka file"})
	}
	defer src.Close()

	xl, err := excelize.OpenReader(src)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Gagal membaca file excel"})
	}

	sheet := xl.GetSheetName(0)
	rows, err := xl.GetRows(sheet)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Gagal membaca sheet"})
	}

	var inserted, updated, skipped int
	var skipReasons = make(map[string]int)

	for i := 1; i < len(rows); i++ {
		if utils.IsBankHeaderRow(rows[i]) {
			skipped++
			skipReasons["header"]++
			continue
		}

		if len(rows[i]) == 0 || utils.IsEmptyRow(rows[i]) {
			skipped++
			skipReasons["empty_row"]++
			continue
		}

		terminalID := dto.GetRawCell(xl, sheet, i, 1)
		namaNasabah := dto.GetRawCell(xl, sheet, i, 5)

		if terminalID == "" || namaNasabah == "" {
			skipped++
			if terminalID == "" {
				skipReasons["no_terminal_id"]++
			} else {
				skipReasons["no_nasabah"]++
			}
			continue
		}

		tanggalPasang, err := dto.GetExcelDate(xl, sheet, i, 6) // ACTUAL_START_DATE (kolom G)
		if err != nil {
			skipped++
			skipReasons["invalid_date"]++
			continue
		}

		var mesin models.MesinEDC
		err = database.DB.
			Where("terminal_id = ?", terminalID).
			First(&mesin).Error

		if err == nil {
			bankData := models.MesinEDC{
				NamaNasabah:   namaNasabah,
				TanggalPasang: tanggalPasang,
			}

			updates := utils.MergeBankData(&mesin, bankData)
			if mesin.LetakMesin == "vendor" {
				updates["letak_mesin"] = "nasabah"
			}

			if len(updates) == 0 {
				skipped++
				skipReasons["no_changes"]++
				continue
			}

			if err := database.DB.Model(&mesin).Updates(updates).Error; err != nil {
				skipped++
				skipReasons["update_failed"]++
				continue
			}
			var sewa models.Sewa
			errSewa := database.DB.Where("mesin_id = ?", mesin.ID).First(&sewa).Error

			if errors.Is(errSewa, gorm.ErrRecordNotFound) {
				sewa = models.Sewa{
					MesinID:      mesin.ID,
					StatusSewa:   "aktif",
					BiayaBulanan: 150000,
				}
				database.DB.Create(&sewa)
			} else if errSewa == nil {
				sewa.StatusSewa = "aktif"
				if sewa.BiayaBulanan <= 0 {
					sewa.BiayaBulanan = 150000
				}
				database.DB.Save(&sewa)
			}

			updated++
			continue
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			newMesin := models.MesinEDC{
				TerminalID:    terminalID,
				NamaNasabah:   namaNasabah,
				TanggalPasang: tanggalPasang,
				StatusData:    "bank",
				StatusMesin:   "aktif",
				LetakMesin:    "nasabah",
			}

			if err := database.DB.Create(&newMesin).Error; err != nil {
				skipped++
				skipReasons["create_failed"]++
				continue
			}

			sewa := models.Sewa{
				MesinID:      newMesin.ID,
				StatusSewa:   "aktif",
				BiayaBulanan: 150000,
			}

			database.DB.Create(&sewa)

			inserted++
			continue
		}

		skipped++
		skipReasons["unknown_error"]++
	}

	return c.JSON(fiber.Map{
		"message":      "Upload bank selesai",
		"inserted":     inserted,
		"updated":      updated,
		"skipped":      skipped,
		"skip_reasons": skipReasons,
		"total_rows":   len(rows) - 1,
	})
}