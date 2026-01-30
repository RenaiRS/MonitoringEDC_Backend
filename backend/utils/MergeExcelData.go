package utils

import "backend/models"

func MergeVendorData(existing *models.MesinEDC, vendor models.MesinEDC) map[string]interface{} {
	updates := map[string]interface{}{}

	if existing.MID == "" && vendor.MID != "" {
		updates["mid"] = vendor.MID
	}

	if existing.Kota == "" && vendor.Kota != "" {
		updates["kota"] = vendor.Kota
	}

	if existing.Cabang == "" && vendor.Cabang != "" {
		updates["cabang"] = vendor.Cabang
	}

	if existing.TipeEDC == "" && vendor.TipeEDC != "" {
		updates["tipe_edc"] = vendor.TipeEDC
	}

	return updates
}

func MergeBankData(existing *models.MesinEDC, bank models.MesinEDC) map[string]interface{} {
	updates := map[string]interface{}{}

	if bank.NamaNasabah != "" {
		updates["nama_nasabah"] = bank.NamaNasabah
	}

	if bank.TanggalPasang != nil {
		updates["tanggal_pasang"] = bank.TanggalPasang
	}

	if existing.StatusData != "bank" {
		updates["status_data"] = "bank"
	}

	return updates
}
