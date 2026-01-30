package controllers

import (
	"time"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"backend/database"
	"backend/models"
	"backend/models/dto"
	"backend/utils"
)

func GetDetailMesin(c *fiber.Ctx) error {
	id := c.Params("id")

	var mesin models.MesinEDC
	if err := database.DB.
		Preload("Perbaikan", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Where("terminal_id = ?", id).
		First(&mesin).Error; err != nil {
		return utils.Error(c, "Mesin tidak ditemukan")
	}

	resp := dto.DetailMesinResponse{}

	resp.InformasiMesin.TerminalID = mesin.TerminalID
	resp.InformasiMesin.MID = mesin.MID
	resp.InformasiMesin.TipeEDC = mesin.TipeEDC
	resp.InformasiMesin.Vendor = mesin.Vendor
	resp.InformasiMesin.StatusMesin = dto.MapStatusMesin(mesin.StatusMesin)
	resp.InformasiMesin.StatusData = dto.MapStatusData(mesin.StatusData)

	resp.InformasiLokasi.NamaNasabah = mesin.NamaNasabah
	if mesin.StatusData == "vendor_only" {
		resp.InformasiLokasi.NamaNasabah = "Belum Terdata"
	}

	resp.InformasiLokasi.Cabang = mesin.Cabang
	resp.InformasiLokasi.Kota = mesin.Kota
	resp.InformasiLokasi.StatusLetak = dto.MapStatusLetak(mesin.LetakMesin)

	if mesin.TanggalPasang != nil {
		resp.InformasiLokasi.TanggalPasang = dto.FormatDateOnlyPtr(mesin.TanggalPasang)
	}
	
	var sewaAktif models.Sewa
	errAktif := database.DB.
		Where("mesin_id = ? AND status_sewa = ?", mesin.ID, "aktif").
		Order("created_at DESC").
		First(&sewaAktif).Error

	if errAktif == nil {
		resp.InformasiSewa.StatusSewa = "AKTIF"
		resp.InformasiSewa.BiayaBulanan = utils.NormalizeBiayaBulanan(sewaAktif.BiayaBulanan)
	} else {
		var sewaLast models.Sewa
		errLast := database.DB.
			Where("mesin_id = ?", mesin.ID).
			Order("created_at DESC").
			First(&sewaLast).Error

		if errLast == nil {
			resp.InformasiSewa.StatusSewa = dto.MapStatusSewa(sewaLast.StatusSewa)
			if sewaLast.StatusSewa == "berakhir" {
				resp.InformasiSewa.BiayaBulanan = 0
			} else {
				resp.InformasiSewa.BiayaBulanan = utils.NormalizeBiayaBulanan(sewaLast.BiayaBulanan)
			}
		} else {
			statusSewa, biayaBulanan := utils.GetBiayaSewaForMesin(database.DB, mesin.ID)
			resp.InformasiSewa.StatusSewa = statusSewa
			resp.InformasiSewa.BiayaBulanan = biayaBulanan
		}
	}

	if len(mesin.Perbaikan) > 0 {
		formatted := dto.FormatDateOnlyPtr(mesin.Perbaikan[0].EstimasiSelesai)
		resp.InformasiSewa.EstimasiSelesai = &formatted
	}

	resp.SumberData = dto.MapStatusData(mesin.StatusData)

	return utils.Success(c, resp)
}

type UpdateMesinRequest struct {
	NamaNasabah     string     `json:"nama_nasabah"`
	Kota            string     `json:"kota"`
	Cabang          string     `json:"cabang"`
	TipeEDC         string     `json:"tipe_edc"`
	Vendor          string     `json:"vendor"`
	StatusMesin     string     `json:"status_mesin"`
	StatusData      string     `json:"status_data"`
	StatusSewa      string     `json:"status_sewa"`
	StatusLetak     string     `json:"status_letak"`
	TanggalPasang   *time.Time `json:"tanggal_pasang"`
	BiayaBulanan    int        `json:"biaya_bulanan"`
	EstimasiSelesai *time.Time `json:"estimasi_selesai"`
}

func UpdateMesin(c *fiber.Ctx) error {
	id := c.Params("id")

	var mesin models.MesinEDC
	if err := database.DB.Where("terminal_id = ?", id).First(&mesin).Error; err != nil {
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
	mesin.StatusMesin = utils.ConvertStatusMesinToDB(req.StatusMesin)
	mesin.StatusData = utils.ConvertStatusDataToDB(req.StatusData)
	mesin.LetakMesin = utils.ConvertStatusLetakToDB(req.StatusLetak)

	if req.TanggalPasang != nil {
		mesin.TanggalPasang = req.TanggalPasang
	}

	if err := database.DB.Save(&mesin).Error; err != nil {
		return utils.Error(c, "Gagal menyimpan data mesin")
	}

	var sewa models.Sewa
	database.DB.Where("mesin_id = ?", mesin.ID).
		FirstOrCreate(&sewa, models.Sewa{MesinID: mesin.ID})

	sewa.StatusSewa = utils.ConvertStatusSewaToDB(req.StatusSewa)
	sewa.BiayaBulanan = req.BiayaBulanan

	if err := database.DB.Save(&sewa).Error; err != nil {
		return utils.Error(c, "Gagal menyimpan data sewa")
	}

	if req.StatusMesin == "PERBAIKAN" && req.EstimasiSelesai != nil {
		var perbaikan models.Perbaikan
		database.DB.Where("mesin_id = ?", mesin.ID).
			Order("created_at DESC").
			FirstOrCreate(&perbaikan, models.Perbaikan{MesinID: mesin.ID})

		perbaikan.StatusPerbaikan = "perbaikan"
		perbaikan.EstimasiSelesai = req.EstimasiSelesai

		if err := database.DB.Save(&perbaikan).Error; err != nil {
			return utils.Error(c, "Gagal menyimpan data perbaikan")
		}
	}

	return utils.Success(c, "Data mesin berhasil diperbarui")
}