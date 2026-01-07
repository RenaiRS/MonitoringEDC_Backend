package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"backend/database"
	"backend/models"
	"backend/utils"
	"backend/models/dto"
)

func GetOverdueSummary(c *fiber.Ctx) error {
	var mesin []models.MesinEDC
	db := database.DB

	err := db.
		Preload("Perbaikan").
		Preload("Sewa").
		Where("status_mesin = ?", "perbaikan").
		Find(&mesin).Error

	if err != nil {
		return utils.Error(c, "Gagal mengambil data overdue")
	}

	now := time.Now()
	totalPerbaikan := 0
	warning := 0
	overdue := 0
	totalKerugian := 0

	for _, m := range mesin {
		if len(m.Perbaikan) == 0 {
			continue
		}

		totalPerbaikan++

		p := m.Perbaikan[0] // 1 mesin = 1 perbaikan aktif

		if p.EstimasiPerbaikan == nil {
			continue
		}

		diff := int(now.Sub(*p.EstimasiPerbaikan).Hours() / 24)

		if diff >= 3 {
			overdue++
			if m.Sewa != nil {
				totalKerugian += m.Sewa.BiayaBulanan
			}
		} else if diff >= 0 {
			warning++
		}
	}

	return utils.Success(c, fiber.Map{
		"total_perbaikan":  totalPerbaikan,
		"warning":          warning,
		"overdue":          overdue,
		"estimasi_kerugian": totalKerugian,
	})
}

func GetOverdueList(c *fiber.Ctx) error {
	var mesin []models.MesinEDC
	db := database.DB

	err := db.
		Preload("Perbaikan").
		Preload("Sewa").
		Where("status_mesin = ?", "perbaikan").
		Find(&mesin).Error

	if err != nil {
		return utils.Error(c, "Gagal mengambil data mesin overdue")
	}

	now := time.Now()
	var result []dto.OverdueMesinResponse

	for _, m := range mesin {
		if len(m.Perbaikan) == 0 {
			continue
		}

		p := m.Perbaikan[0]

		status := "perbaikan"
		terlambat := 0
		kerugian := 0

		if p.EstimasiPerbaikan != nil {
			diff := int(now.Sub(*p.EstimasiPerbaikan).Hours() / 24)

			if diff >= 3 {
				status = "overdue"
				terlambat = diff
				if m.Sewa != nil {
					kerugian = m.Sewa.BiayaBulanan
				}
			} else if diff >= 0 {
				status = "warning"
				terlambat = diff
			}
		}

		result = append(result, dto.OverdueMesinResponse{
			TerminalID:        m.TerminalID,
			NamaNasabah:       utils.GetNamaNasabah(m.NamaNasabah, m.StatusData),
			LokasiMesin:       m.LetakMesin,
			TanggalPasang:     m.TanggalPasang,
			EstimasiPerbaikan: p.EstimasiPerbaikan,
			TerlambatHari:     terlambat,
			StatusPerbaikan:   status,
			Kerugian:          kerugian,
		})
	}

	return utils.Success(c, result)
}

func SearchOverdue(c *fiber.Ctx) error {
	query := c.Query("q")
	db := database.DB

	var mesin []models.MesinEDC
	err := db.
		Preload("Perbaikan").
		Preload("Sewa").
		Where(`
			status_mesin = ?
			AND (terminal_id LIKE ? OR nama_nasabah LIKE ?)
		`, "perbaikan", "%"+query+"%", "%"+query+"%").
		Find(&mesin).Error

	if err != nil {
		return utils.Error(c, "Gagal melakukan pencarian overdue")
	}

	now := time.Now()
	var result []dto.OverdueMesinResponse

	for _, m := range mesin {
		if len(m.Perbaikan) == 0 {
			continue
		}

		p := m.Perbaikan[0]

		status := "perbaikan"
		terlambat := 0
		kerugian := 0

		if p.EstimasiPerbaikan != nil {
			diff := int(now.Sub(*p.EstimasiPerbaikan).Hours() / 24)

			if diff >= 3 {
				status = "overdue"
				terlambat = diff
				if m.Sewa != nil {
					kerugian = m.Sewa.BiayaBulanan
				}
			} else if diff >= 0 {
				status = "warning"
				terlambat = diff
			}
		}

		result = append(result, dto.OverdueMesinResponse{
			TerminalID:        m.TerminalID,
			NamaNasabah:       utils.GetNamaNasabah(m.NamaNasabah, m.StatusData),
			LokasiMesin:       m.LetakMesin,
			TanggalPasang:     m.TanggalPasang,
			EstimasiPerbaikan: p.EstimasiPerbaikan,
			TerlambatHari:     terlambat,
			StatusPerbaikan:   status,
			Kerugian:          kerugian,
		})
	}

	return utils.Success(c, result)
}
