package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	// "gorm.io/gorm"

	"backend/database"
	"backend/models"
	"backend/utils"
)

func GetTotalMesin(c *fiber.Ctx) error {
	var total int64

	database.DB.Model(&models.MesinEDC{}).Count(&total)

	return c.JSON(fiber.Map{
		"total_mesin": total,
	})
}

func GetMesinTerdataBank(c *fiber.Ctx) error {
	var total int64

	database.DB.
		Model(&models.MesinEDC{}).
		Where("status_data = ?", "bank").
		Count(&total)

	return c.JSON(fiber.Map{
		"mesin_terdata_bank": total,
	})
}

// func GetMesinVendorOnly(c *fiber.Ctx) error {
// 	var mesin []models.MesinEDC

// 	database.DB.
// 		Where("status_data = ?", "vendor_only").
// 		Order("created_at DESC").
// 		Find(&mesin)

// 	return c.JSON(fiber.Map{
// 		"mesin_vendor_only": mesin,
// 	})
// }

func GetOverdueStatus(c *fiber.Ctx) error {
	type Result struct {
		Status string `json:"status"`
		Total  int64  `json:"total"`
	}

	var results []Result

	err := database.DB.
		Model(&models.Perbaikan{}).
		Select("status_perbaikan as status, COUNT(*) as total").
		Group("status_perbaikan").
		Scan(&results).Error

	if err != nil {
		return utils.Error(c, "Gagal mengambil status overdue")
	}

	return utils.Success(c, results)
}

func GetStatusMesin(c *fiber.Ctx) error {
	type Result struct {
		StatusMesin string
		Total       int64
	}

	var result []Result

	database.DB.
		Model(&models.MesinEDC{}).
		Select("status_mesin, COUNT(*) as total").
		Group("status_mesin").
		Scan(&result)

	return c.JSON(result)
}

func GetMesinBaruVendor(c *fiber.Ctx) error {
	var mesin []models.MesinEDC

	err := database.DB.
		Where("status_data = ?", "vendor_only").
		Order("created_at DESC").
		Find(&mesin).Error

	if err != nil {
		return utils.Error(c, "Gagal mengambil mesin baru dari vendor")
	}

	return utils.Success(c, mesin)
}

func GetMonitoringOverdue(c *fiber.Ctx) error {
	now := time.Now()

	var data []models.Perbaikan

	database.DB.
		Preload("Mesin").
		Where("estimasi_selesai < ? AND status_perbaikan = ?", now, "overdue").
		Find(&data)

	return c.JSON(fiber.Map{
		"monitoring_overdue": data,
	})
}

// func GetTotalOverdue(c *fiber.Ctx) error {
// 	now := time.Now()
// 	var total int64

// 	database.DB.
// 		Model(&models.Perbaikan{}).
// 		Where("estimasi_selesai < ? AND status_perbaikan = ?", now, "overdue").
// 		Count(&total)

// 	return c.JSON(fiber.Map{
// 		"total_overdue": total,
// 	})
// }

