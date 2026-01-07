package controllers

import (
	"github.com/gofiber/fiber/v2"
	"backend/database"
	"backend/models"
)

func GetSewaSummary(c *fiber.Ctx) error {
	db := database.DB

	var sewas []models.Sewa
	db.Preload("Mesin").Find(&sewas)

	sewaAktif := 0
	sewaBerakhir := 0
	bermasalah := 0
	totalBiaya := 0

	for _, s := range sewas {
		biaya := s.BiayaBulanan
		if biaya == 0 {
			biaya = 1500000
		}

		if s.StatusSewa == "aktif" {
			sewaAktif++
			totalBiaya += biaya
		} else if s.StatusSewa == "berakhir" {
			sewaBerakhir++
		}

		if s.Mesin != nil {
			if s.Mesin.StatusMesin == "perbaikan" || s.Mesin.StatusMesin == "rusak" {
				bermasalah++
			}
		}
	}

	return c.JSON(fiber.Map{
		"sewa_aktif":          sewaAktif,
		"sewa_berakhir":       sewaBerakhir,
		"total_biaya_bulanan": totalBiaya,
		"bermasalah":          bermasalah,
	})
}

func GetSewaList(c *fiber.Ctx) error {
	db := database.DB

	var sewas []models.Sewa
	err := db.Preload("Mesin").Find(&sewas).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal mengambil data sewa",
		})
	}

	var result []fiber.Map

	for _, s := range sewas {
		if s.Mesin == nil {
			continue
		}

		biaya := s.BiayaBulanan
		if biaya == 0 {
			biaya = 1500000
		}

		result = append(result, fiber.Map{
			"terminal_id":     s.Mesin.TerminalID,
			"nama_nasabah":    s.Mesin.NamaNasabah,
			"tanggal_pasang":  s.Mesin.TanggalPasang,
			"status_mesin":    s.Mesin.StatusMesin,
			"letak_mesin":     s.Mesin.LetakMesin,
			"status_sewa":     s.StatusSewa,
			"biaya_bulanan":   biaya,
		})
	}

	return c.JSON(result)
}

func SearchSewa(c *fiber.Ctx) error {
	query := c.Query("q")
	db := database.DB

	var sewas []models.Sewa
	err := db.
		Joins("JOIN mesin_edcs ON mesin_edcs.id = sewas.mesin_id").
		Where(`
			mesin_edcs.terminal_id LIKE ? OR mesin_edcs.nama_nasabah LIKE ?
		`, "%"+query+"%", "%"+query+"%").
		Preload("Mesin").
		Find(&sewas).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal melakukan pencarian sewa",
		})
	}

	var result []fiber.Map

	for _, s := range sewas {
		if s.Mesin == nil {
			continue
		}

		biaya := s.BiayaBulanan
		if biaya == 0 {
			biaya = 1500000
		}

		result = append(result, fiber.Map{
			"terminal_id":     s.Mesin.TerminalID,
			"nama_nasabah":    s.Mesin.NamaNasabah,
			"tanggal_pasang":  s.Mesin.TanggalPasang,
			"status_mesin":    s.Mesin.StatusMesin,
			"letak_mesin":     s.Mesin.LetakMesin,
			"status_sewa":     s.StatusSewa,
			"biaya_bulanan":   biaya,
		})
	}

	return c.JSON(result)
}

