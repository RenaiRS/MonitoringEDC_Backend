package controllers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"backend/database"
	"backend/models"
	"backend/utils"
	"backend/models/dto"
)

func GetRekapMesin(c *fiber.Ctx) error {
	search := c.Query("search")

	var mesin []models.MesinEDC

	query := database.DB.Model(&models.MesinEDC{})

	if search != "" {
		keyword := "%" + strings.ToLower(search) + "%"
		query = query.Where(
			"LOWER(terminal_id) LIKE ? OR LOWER(nama_nasabah) LIKE ?",
			keyword, keyword,
		)
	}

	if err := query.Order("tanggal_pasang DESC").Find(&mesin).Error; err != nil {
		return utils.Error(c, "Gagal mengambil data rekap")
	}

	var response []dto.RekapMesinResponse
	for _, m := range mesin {
		response = append(response, dto.RekapMesinResponse{
			ID:           m.ID,
			TerminalID:   m.TerminalID,
			NamaNasabah:  m.NamaNasabah,
			Cabang:       m.Cabang,
			TanggalPasang: m.TanggalPasang,
			StatusMesin:  m.StatusMesin,
			StatusData:   m.StatusData,
		})
	}

	return utils.Success(c, response)
}

func CreateRekapMesin(c *fiber.Ctx) error {
	type Request struct {
		TerminalID        string    `json:"terminal_id"`
		MID               string    `json:"mid"`
		Kota              string    `json:"kota"`
		Cabang            string    `json:"cabang"`
		TipeEDC           string    `json:"tipe_edc"`
		TanggalPasang     time.Time `json:"tanggal_pasang"`
		BiayaSewaBulanan  int       `json:"biaya_sewa_bulanan"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, "Request tidak valid")
	}

	// ✅ VALIDASI WAJIB
	if req.TerminalID == "" || req.MID == "" {
		return utils.Error(c, "Terminal ID dan MID wajib diisi")
	}

	// =========================
	// CREATE MESIN (DEFAULT)
	// =========================
	mesin := models.MesinEDC{
		TerminalID:    req.TerminalID,
		MID:           req.MID,
		Kota:          req.Kota,
		Cabang:        req.Cabang,
		TipeEDC:       req.TipeEDC,
		TanggalPasang: req.TanggalPasang,

		StatusData:   "vendor_only",
		StatusMesin:  "aktif",
		LetakMesin:   "nasabah",
	}

	if err := database.DB.Create(&mesin).Error; err != nil {
		return utils.Error(c, "Gagal menyimpan mesin")
	}

	// =========================
	// CREATE SEWA (DEFAULT)
	// =========================
	biaya := req.BiayaSewaBulanan
	if biaya == 0 {
		biaya = 150000 // default
	}

	sewa := models.Sewa{
		MesinID:      mesin.ID,
		BiayaBulanan: biaya,
		StatusSewa:   "berakhir",
	}

	if err := database.DB.Create(&sewa).Error; err != nil {
		return utils.Error(c, "Gagal menyimpan data sewa")
	}

	return utils.Success(c, fiber.Map{
		"message": "Rekap mesin berhasil ditambahkan",
		"id_mesin": mesin.ID,
	})
}
