package controllers

import (
	"time"

	"backend/database"
	"backend/models"
	"backend/models/dto"
	"backend/utils"

	"github.com/gofiber/fiber/v2"
)

func GetOverdueSummary(c *fiber.Ctx) error {
	db := database.DB

	var mesin []models.MesinEDC
	err := db.Preload("Perbaikan").
		Where("status_mesin = ?", "perbaikan").
		Find(&mesin).Error

	if err != nil {
		return utils.Error(c, "Gagal mengambil data overdue")
	}

	var sewasAktif []models.Sewa
	err = db.Where("status_sewa = ?", "aktif").
		Order("created_at DESC").
		Find(&sewasAktif).Error

	if err != nil {
		return utils.Error(c, "Gagal mengambil data sewa aktif")
	}

	sewaMap := make(map[uint]*models.Sewa)
	for i := range sewasAktif {
		mesinID := sewasAktif[i].MesinID
		if _, exists := sewaMap[mesinID]; !exists {
			sewaMap[mesinID] = &sewasAktif[i]
		}
	}

	totalPerbaikan := 0
	warning := 0
	overdue := 0
	totalKerugian := 0

	now := time.Now()

	for _, m := range mesin {
		if len(m.Perbaikan) == 0 || m.Perbaikan[0].EstimasiSelesai == nil {
			continue
		}

		totalPerbaikan++
		p := m.Perbaikan[0]

		diffDays := int(p.EstimasiSelesai.Sub(now).Hours() / 24)

		if diffDays < 0 {
			overdue++
			daysLate := -diffDays

			if sewa, exists := sewaMap[m.ID]; exists {
				biayaBulanan := utils.NormalizeBiayaBulanan(sewa.BiayaBulanan)
				
				bulanOverdue := (daysLate / 30) + 1
				kerugian := biayaBulanan * bulanOverdue
				
				totalKerugian += kerugian
			}

		} else if diffDays <= 3 {
			warning++
		}
	}

	statusOverdue := []fiber.Map{
		{"status": "PERBAIKAN", "total": totalPerbaikan - warning - overdue},
		{"status": "WARNING", "total": warning},
		{"status": "OVERDUE", "total": overdue},
	}

	return utils.Success(c, fiber.Map{
		"total_perbaikan":   totalPerbaikan,
		"warning":           warning,
		"overdue":           overdue,
		"statusOverdue":     statusOverdue,
		"estimasi_kerugian": totalKerugian,
	})
}

func GetOverdueList(c *fiber.Ctx) error {
	db := database.DB

	var mesin []models.MesinEDC
	err := db.Preload("Perbaikan").
		Where("status_mesin = ?", "perbaikan").
		Find(&mesin).Error

	if err != nil {
		return utils.Error(c, "Gagal mengambil data mesin overdue")
	}

	var sewasAktif []models.Sewa
	err = db.Where("status_sewa = ?", "aktif").
		Order("created_at DESC").
		Find(&sewasAktif).Error

	if err != nil {
		return utils.Error(c, "Gagal mengambil data sewa aktif")
	}

	sewaMap := make(map[uint]*models.Sewa)
	for i := range sewasAktif {
		mesinID := sewasAktif[i].MesinID
		if _, exists := sewaMap[mesinID]; !exists {
			sewaMap[mesinID] = &sewasAktif[i]
		}
	}

	now := time.Now()
	var result []dto.MachineResponse

	for _, m := range mesin {
		if len(m.Perbaikan) == 0 || m.Perbaikan[0].EstimasiSelesai == nil {
			continue
		}

		p := m.Perbaikan[0]
		status := "PERBAIKAN"
		kerugian := 0
		daysLate := 0

		diffDays := int(p.EstimasiSelesai.Sub(now).Hours() / 24)

		if diffDays < 0 {
			status = "OVERDUE"
			daysLate = -diffDays

			if sewa, exists := sewaMap[m.ID]; exists {
				biayaBulanan := utils.NormalizeBiayaBulanan(sewa.BiayaBulanan)
				
				bulanOverdue := (daysLate / 30) + 1
				kerugian = biayaBulanan * bulanOverdue
			}

		} else if diffDays <= 3 {
			status = "WARNING"
		}

		resp := dto.MachineResponse{
			ID:              m.ID,
			TerminalID:      m.TerminalID,
			MID:             m.MID,
			NamaNasabah:     utils.GetNamaNasabah(m.NamaNasabah, m.StatusData),
			Kota:            m.Kota,
			Cabang:          m.Cabang,
			TipeEDC:         m.TipeEDC,
			StatusMesin:     status,
			StatusData:      dto.MapStatusData(m.StatusData),
			StatusLetak:     dto.MapStatusLetak(m.LetakMesin),
			BiayaSewa:       kerugian,
			DaysOverdue:     daysLate,
			StatusPerbaikan: status,
		}

		if tp := dto.FormatDateOnlyPtr(m.TanggalPasang); tp != "" {
			resp.TanggalPasang = tp
		}
		if es := dto.FormatDateOnlyPtr(p.EstimasiSelesai); es != "" {
			resp.EstimasiSelesai = &es
		}

		result = append(result, resp)
	}

	return utils.Success(c, result)
}

func SearchOverdue(c *fiber.Ctx) error {
	query := c.Query("q")
	db := database.DB

	var mesin []models.MesinEDC
	err := db.Preload("Perbaikan").
		Where("status_mesin = ? AND (terminal_id LIKE ? OR nama_nasabah LIKE ?)",
			"perbaikan", "%"+query+"%", "%"+query+"%").
		Find(&mesin).Error

	if err != nil {
		return utils.Error(c, "Gagal melakukan pencarian overdue")
	}

	var sewasAktif []models.Sewa
	err = db.Where("status_sewa = ?", "aktif").
		Order("created_at DESC").
		Find(&sewasAktif).Error

	if err != nil {
		return utils.Error(c, "Gagal mengambil data sewa aktif")
	}

	sewaMap := make(map[uint]*models.Sewa)
	for i := range sewasAktif {
		mesinID := sewasAktif[i].MesinID
		if _, exists := sewaMap[mesinID]; !exists {
			sewaMap[mesinID] = &sewasAktif[i]
		}
	}

	now := time.Now()
	var result []dto.MachineResponse

	for _, m := range mesin {
		if len(m.Perbaikan) == 0 || m.Perbaikan[0].EstimasiSelesai == nil {
			continue
		}

		p := m.Perbaikan[0]
		status := "PERBAIKAN"
		kerugian := 0
		daysLate := 0

		diffDays := int(p.EstimasiSelesai.Sub(now).Hours() / 24)

		if diffDays < 0 {
			status = "OVERDUE"
			daysLate = -diffDays

			if sewa, exists := sewaMap[m.ID]; exists {
				biayaBulanan := utils.NormalizeBiayaBulanan(sewa.BiayaBulanan)
				
				bulanOverdue := (daysLate / 30) + 1
				kerugian = biayaBulanan * bulanOverdue
			}
		} else if diffDays <= 3 {
			status = "WARNING"
		}

		resp := dto.MachineResponse{
			ID:              m.ID,
			TerminalID:      m.TerminalID,
			MID:             m.MID,
			NamaNasabah:     utils.GetNamaNasabah(m.NamaNasabah, m.StatusData),
			Kota:            m.Kota,
			Cabang:          m.Cabang,
			TipeEDC:         m.TipeEDC,
			StatusMesin:     status,
			StatusData:      dto.MapStatusData(m.StatusData),
			StatusLetak:     dto.MapStatusLetak(m.LetakMesin),
			BiayaSewa:       kerugian,
			DaysOverdue:     daysLate,
			StatusPerbaikan: status,
		}

		if tp := dto.FormatDateOnlyPtr(m.TanggalPasang); tp != "" {
			resp.TanggalPasang = tp
		}
		if es := dto.FormatDateOnlyPtr(p.EstimasiSelesai); es != "" {
			resp.EstimasiSelesai = &es
		}

		result = append(result, resp)
	}

	return utils.Success(c, result)
}