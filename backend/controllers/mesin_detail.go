package controllers

import (
	"github.com/gofiber/fiber/v2"
	"time"
	"backend/database"
	"backend/models"
	"backend/models/dto"
	"backend/utils"
)

func GetDetailMesin(c *fiber.Ctx) error {
	id := c.Params("id")

	var mesin models.MesinEDC
	if err := database.DB.First(&mesin, id).Error; err != nil {
		return utils.Error(c, "Mesin tidak ditemukan")
	}

	var sewa models.Sewa
	database.DB.Where("mesin_id = ?", mesin.ID).First(&sewa)

	var perbaikan models.Perbaikan
	database.DB.Where("mesin_id = ?", mesin.ID).
		Order("created_at DESC").
		First(&perbaikan)

	resp := dto.DetailMesinResponse{}

	resp.InformasiMesin.TerminalID = mesin.TerminalID
	resp.InformasiMesin.MID = mesin.MID
	resp.InformasiMesin.TipeEDC = mesin.TipeEDC
	resp.InformasiMesin.Vendor = mesin.Vendor
	resp.InformasiMesin.StatusMesin = mesin.StatusMesin
	resp.InformasiMesin.StatusData = mesin.StatusData

	resp.InformasiLokasi.NamaNasabah = mesin.NamaNasabah
	if mesin.StatusData == "vendor_only" {
		resp.InformasiLokasi.NamaNasabah = "Belum Terdata"
	}
	resp.InformasiLokasi.Cabang = mesin.Cabang
	resp.InformasiLokasi.Kota = mesin.Kota
	resp.InformasiLokasi.TanggalPasang = mesin.TanggalPasang

	resp.InformasiSewa.StatusSewa = sewa.StatusSewa
	resp.InformasiSewa.BiayaBulanan = sewa.BiayaBulanan
	resp.InformasiSewa.EstimasiSelesai = perbaikan.EstimasiSelesai

	resp.SumberData = mesin.StatusData

	return utils.Success(c, resp)
}

type UpdateMesinRequest struct {
	NamaNasabah    string     `json:"nama_nasabah"`
	Kota           string     `json:"kota"`
	Cabang         string     `json:"cabang"`
	TipeEDC        string     `json:"tipe_edc"`
	Vendor         string     `json:"vendor"`
	StatusMesin    string     `json:"status_mesin"`
	StatusData     string     `json:"status_data"`
	StatusSewa     string     `json:"status_sewa"`
	TanggalPasang  time.Time  `json:"tanggal_pasang"`
	BiayaBulanan   int        `json:"biaya_bulanan"`
	EstimasiSelesai *time.Time `json:"estimasi_selesai"`
}

func UpdateMesin(c *fiber.Ctx) error {
	id := c.Params("id")

	var mesin models.MesinEDC
	if err := database.DB.First(&mesin, id).Error; err != nil {
		return utils.Error(c, "Mesin tidak ditemukan")
	}

	var req UpdateMesinRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, "Input tidak valid")
	}

	mesin.NamaNasabah = req.NamaNasabah
	mesin.Kota = req.Kota
	mesin.Cabang = req.Cabang
	mesin.TipeEDC = req.TipeEDC
	mesin.Vendor = req.Vendor
	mesin.StatusMesin = req.StatusMesin
	mesin.StatusData = req.StatusData
	mesin.TanggalPasang = req.TanggalPasang

	database.DB.Save(&mesin)

	var sewa models.Sewa
	database.DB.Where("mesin_id = ?", mesin.ID).FirstOrCreate(&sewa)
	sewa.StatusSewa = req.StatusSewa
	sewa.BiayaBulanan = req.BiayaBulanan
	database.DB.Save(&sewa)

	if req.StatusMesin == "perbaikan" {
		var perbaikan models.Perbaikan
		database.DB.Where("mesin_id = ?", mesin.ID).
			FirstOrCreate(&perbaikan)
		perbaikan.StatusPerbaikan = "perbaikan"
		perbaikan.EstimasiSelesai = req.EstimasiSelesai
		database.DB.Save(&perbaikan)
	}

	return utils.Success(c, "Data mesin berhasil diperbarui")
}
