package controllers

import (
	"strings"
	"time"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"backend/database"
	"backend/models"
	"backend/models/dto"
	"backend/utils"
)

func GetRekapMesin(c *fiber.Ctx) error {
	search := c.Query("search")

	var mesinList []models.MesinEDC

	query := database.DB.
		Preload("Sewas", func(db *gorm.DB) *gorm.DB {
			return db.
				Where("status_sewa = ?", "aktif").
				Order("created_at DESC").
				Limit(1)
		}).
		Preload("Perbaikan")

	if search != "" {
		keyword := "%" + strings.ToLower(search) + "%"
		query = query.Where(
			"LOWER(terminal_id) LIKE ? OR LOWER(nama_nasabah) LIKE ?",
			keyword, keyword,
		)
	}

	if err := query.Order("tanggal_pasang DESC").Find(&mesinList).Error; err != nil {
		return utils.Error(c, "Gagal mengambil data rekap")
	}

	var response []dto.MachineResponse

	for _, m := range mesinList {
		resp := dto.MachineResponse{
			ID:          m.ID,
			TerminalID:  m.TerminalID,
			MID:         m.MID,
			NamaNasabah: m.NamaNasabah,
			Kota:        m.Kota,
			Cabang:      m.Cabang,
			TipeEDC:     m.TipeEDC,
			Vendor:      m.Vendor,

			StatusMesin: dto.MapStatusMesin(m.StatusMesin),
			StatusData:  dto.MapStatusData(m.StatusData),
			StatusLetak: dto.MapStatusLetak(m.LetakMesin),

			CreatedAt: dto.FormatDateOnly(m.CreatedAt),
			UpdatedAt: dto.FormatDateOnly(m.UpdatedAt),
		}

		if tp := dto.FormatDateOnlyPtr(m.TanggalPasang); tp != "" {
			resp.TanggalPasang = tp
		}

		if len(m.Sewas) > 0 {
			sewa := m.Sewas[0]
			resp.StatusSewa = dto.MapStatusSewa(sewa.StatusSewa)
			resp.BiayaSewa = utils.NormalizeBiayaBulanan(sewa.BiayaBulanan)
		} else {
			resp.StatusSewa = "BERAKHIR"
			resp.BiayaSewa = 0
		}

		if len(m.Perbaikan) > 0 {
			es := dto.FormatDateOnlyPtr(m.Perbaikan[0].EstimasiSelesai)
			if es != "" {
				resp.EstimasiSelesai = &es
			}
		}

		response = append(response, resp)
	}

	return utils.Success(c, response)
}

func CreateRekapMesin(c *fiber.Ctx) error {
	type Request struct {
		TerminalID    string     `json:"terminal_id"`
		MID           string     `json:"mid"`
		NamaNasabah   *string    `json:"nama_nasabah"`
		Kota          string     `json:"kota"`
		Cabang        string     `json:"cabang"`
		TipeEDC       string     `json:"tipe_edc"`
		StatusData    string     `json:"status_data"`
		StatusMesin   string     `json:"status_mesin"`
		StatusSewa    string     `json:"status_sewa"`
		StatusLetak   string     `json:"status_letak"`
		TanggalPasang *time.Time `json:"tanggal_pasang"`
		BiayaSewa     int        `json:"biaya_sewa"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		println("Body parse error:", err.Error())
		return utils.Error(c, "Request tidak valid")
	}

	if req.TerminalID == "" || req.MID == "" {
		return utils.Error(c, "Terminal ID dan MID wajib diisi")
	}

	println("Received data:", req.TerminalID, req.MID)

	var existing models.MesinEDC
	if err := database.DB.
		Where("terminal_id = ?", req.TerminalID).
		First(&existing).Error; err == nil {
		return utils.Error(c, "Terminal ID sudah terdaftar")
	}

	namaNasabah := "Belum Terdata"
	if req.NamaNasabah != nil && *req.NamaNasabah != "" {
		namaNasabah = *req.NamaNasabah
	}

	statusData := utils.ConvertStatusDataToDB(req.StatusData)
	statusMesin := utils.ConvertStatusMesinToDB(req.StatusMesin)
	statusLetak := utils.ConvertStatusLetakToDB(req.StatusLetak)

	mesin := models.MesinEDC{
		TerminalID:    req.TerminalID,
		MID:           req.MID,
		NamaNasabah:   namaNasabah,
		Kota:          req.Kota,
		Cabang:        req.Cabang,
		TipeEDC:       req.TipeEDC,
		TanggalPasang: req.TanggalPasang,
		StatusData:    statusData,
		StatusMesin:   statusMesin,
		LetakMesin:    statusLetak,
	}

	if err := database.DB.Create(&mesin).Error; err != nil {
		return utils.Error(c, "Gagal menyimpan mesin: "+err.Error())
	}

	biaya := utils.NormalizeBiayaBulanan(req.BiayaSewa)
	statusSewa := utils.ConvertStatusSewaToDB(req.StatusSewa)

	sewa := models.Sewa{
		MesinID:      mesin.ID,
		BiayaBulanan: biaya,
		StatusSewa:   statusSewa,
	}

	if err := database.DB.Create(&sewa).Error; err != nil {
		return utils.Error(c, "Gagal menyimpan data sewa: "+err.Error())
	}

	return utils.Success(c, fiber.Map{
		"message":     "Rekap mesin berhasil ditambahkan",
		"id_mesin":    mesin.ID,
		"terminal_id": mesin.TerminalID,
	})
}