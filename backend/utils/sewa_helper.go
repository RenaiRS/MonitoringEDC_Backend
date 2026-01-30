package utils

import (
	"backend/models"
	"gorm.io/gorm"
)

func NormalizeBiayaBulanan(biaya int) int {
	if biaya <= 0 {
		return 150000
	}
	return biaya
}

func GetBiayaSewaForMesin(db *gorm.DB, mesinID uint) (statusSewa string, biaya int) {
	var sewaAktif models.Sewa
	errAktif := db.
		Where("mesin_id = ? AND status_sewa = ?", mesinID, "aktif").
		Order("created_at DESC").
		First(&sewaAktif).Error

	if errAktif == nil {
		statusSewa = "AKTIF"
		biaya = NormalizeBiayaBulanan(sewaAktif.BiayaBulanan)
		return
	}

	var sewaLast models.Sewa
	errLast := db.
		Where("mesin_id = ?", mesinID).
		Order("created_at DESC").
		First(&sewaLast).Error

	if errLast == nil {
		statusSewa = MapStatusSewaToDTO(sewaLast.StatusSewa)
		if sewaLast.StatusSewa == "berakhir" {
			biaya = 0
		} else {
			biaya = NormalizeBiayaBulanan(sewaLast.BiayaBulanan)
		}
		return
	}

	statusSewa = "BERAKHIR"
	biaya = 0
	return
}

// IsMesinBermasalah checks if mesin has problematic status
func IsMesinBermasalah(status string) bool {
	return status == "perbaikan" || status == "rusak"
}