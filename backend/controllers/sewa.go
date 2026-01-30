package controllers

import (
	"github.com/gofiber/fiber/v2"
	"backend/database"
	"backend/models"
	"backend/utils"
	"backend/models/dto"
)

func GetSewaSummary(c *fiber.Ctx) error {
	db := database.DB

	var mesins []models.MesinEDC
	err := db.Find(&mesins).Error
	if err != nil {
		return utils.Error(c, "Gagal mengambil summary sewa")
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

	sewaAktif := 0
	sewaBerakhir := 0
	bermasalah := 0
	totalBiaya := 0

	for _, m := range mesins {
		if sewa, exists := sewaMap[m.ID]; exists {
			// Mesin ini punya sewa aktif
			sewaAktif++
			totalBiaya += utils.NormalizeBiayaBulanan(sewa.BiayaBulanan)

			if utils.IsMesinBermasalah(m.StatusMesin) {
				bermasalah++
			}
		} else {
			sewaBerakhir++
		}
	}

	return utils.Success(c, fiber.Map{
		"sewa_aktif":          sewaAktif,
		"sewa_berakhir":       sewaBerakhir,
		"total_biaya_bulanan": totalBiaya,
		"bermasalah":          bermasalah,
	})
}

func GetSewaList(c *fiber.Ctx) error {
	db := database.DB

	var mesins []models.MesinEDC
	err := db.Find(&mesins).Error
	if err != nil {
		return utils.Error(c, "Gagal mengambil data sewa")
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

	var result []dto.MachineResponse

	for _, m := range mesins {
		statusSewa := "BERAKHIR"
		biaya := 0

		if sewa, exists := sewaMap[m.ID]; exists {
			statusSewa = utils.MapStatusSewaToDTO(sewa.StatusSewa)
			biaya = utils.NormalizeBiayaBulanan(sewa.BiayaBulanan)
		}

		resp := dto.MachineResponse{
			ID:          m.ID,
			TerminalID:  m.TerminalID,
			MID:         m.MID,
			NamaNasabah: utils.GetNamaNasabah(m.NamaNasabah, m.StatusData),
			Kota:        m.Kota,
			Cabang:      m.Cabang,
			TipeEDC:     m.TipeEDC,
			StatusMesin: dto.MapStatusMesin(m.StatusMesin),
			StatusData:  dto.MapStatusData(m.StatusData),
			StatusSewa:  statusSewa,
			StatusLetak: dto.MapStatusLetak(m.LetakMesin),
			BiayaSewa:   biaya,
		}

		if tp := dto.FormatDateOnlyPtr(m.TanggalPasang); tp != "" {
			resp.TanggalPasang = tp
		}

		result = append(result, resp)
	}

	return utils.Success(c, result)
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
		return utils.Error(c, "Gagal melakukan pencarian sewa")
	}

	var result []dto.MachineResponse

	for _, s := range sewas {
		if s.Mesin == nil {
			continue
		}

		m := s.Mesin
		biaya := utils.NormalizeBiayaBulanan(s.BiayaBulanan)

		machineResp := dto.MachineResponse{
			ID:          m.ID,
			TerminalID:  m.TerminalID,
			MID:         m.MID,
			NamaNasabah: utils.GetNamaNasabah(m.NamaNasabah, m.StatusData),
			Kota:        m.Kota,
			Cabang:      m.Cabang,
			TipeEDC:     m.TipeEDC,
			StatusMesin: dto.MapStatusMesin(m.StatusMesin),
			StatusData:  dto.MapStatusData(m.StatusData),
			StatusSewa:  utils.MapStatusSewaToDTO(s.StatusSewa),
			StatusLetak: dto.MapStatusLetak(m.LetakMesin),
			BiayaSewa:   biaya,
		}

		tp := dto.FormatDateOnlyPtr(m.TanggalPasang)
		if tp != "" {
			machineResp.TanggalPasang = tp
		}

		result = append(result, machineResp)
	}

	return utils.Success(c, result)
}