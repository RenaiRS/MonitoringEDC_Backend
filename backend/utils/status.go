package utils

const (
	// Status Mesin
	MESIN_AKTIF       = "aktif"
	MESIN_PERBAIKAN   = "perbaikan"
	MESIN_RUSAK       = "rusak"
	MESIN_TIDAK_AKTIF = "tidak_aktif"

	// Status Data
	DATA_BANK       = "bank"
	DATA_VENDORONLY = "vendor_only"

	// Status Sewa
	SEWA_AKTIF        = "aktif"
	SEWA_BELUM_TERDATA = "belum_terdata"
	SEWA_BERAKHIR     = "berakhir"

	// Status Perbaikan / Overdue
	PERBAIKAN = "perbaikan"
	WARNING   = "warning"
	OVERDUE   = "overdue"

	// Letak Mesin
	LETAK_NASABAH = "nasabah"
	LETAK_BANK   = "bank"
	LETAK_VENDOR = "vendor"
)

func MapStatusMesinToDTO(status string) string {
	statusMap := map[string]string{
		"aktif":       "AKTIF",
		"perbaikan":   "PERBAIKAN",
		"rusak":       "RUSAK",
		"tidak_aktif": "NONAKTIF",
		"nonaktif":    "NONAKTIF",
	}
	if mapped, ok := statusMap[status]; ok {
		return mapped
	}
	return "AKTIF"
}

func MapStatusDataToDTO(status string) string {
	statusMap := map[string]string{
		"bank":        "TERDATA_BANK",
		"vendor_only": "VENDOR_ONLY",
	}
	if mapped, ok := statusMap[status]; ok {
		return mapped
	}
	return "VENDOR_ONLY"
}

func MapStatusSewaToDTO(status string) string {
	statusMap := map[string]string{
		"aktif":    "AKTIF",
		"berakhir": "BERAKHIR",
		"AKTIF":    "AKTIF",
		"BERAKHIR": "BERAKHIR",
	}
	if mapped, ok := statusMap[status]; ok {
		return mapped
	}
	return "BERAKHIR"
}

func MapStatusLetakToDTO(letak string) string {
	letakMap := map[string]string{
		"nasabah": "NASABAH",
		"vendor":  "VENDOR",
		"bank":    "BANK",
	}
	if mapped, ok := letakMap[letak]; ok {
		return mapped
	}
	return "NASABAH"
}

func ConvertStatusMesinToDB(status string) string {
	mapper := map[string]string{
		"AKTIF":     "aktif",
		"PERBAIKAN": "perbaikan",
		"RUSAK":     "rusak",
		"NONAKTIF":  "tidak_aktif",
	}
	if v, ok := mapper[status]; ok {
		return v
	}
	return "aktif"
}

// ConvertStatusDataToDB converts frontend status data to database format
func ConvertStatusDataToDB(status string) string {
	mapper := map[string]string{
		"TERDATA_BANK": "bank",
		"VENDOR_ONLY":  "vendor_only",
	}
	if v, ok := mapper[status]; ok {
		return v
	}
	return "vendor_only"
}

// ConvertStatusSewaToDB converts frontend status sewa to database format
func ConvertStatusSewaToDB(status string) string {
	mapper := map[string]string{
		"AKTIF":    "aktif",
		"BERAKHIR": "berakhir",
	}
	if v, ok := mapper[status]; ok {
		return v
	}
	return "berakhir"
}

// ConvertStatusLetakToDB converts frontend status letak to database format
func ConvertStatusLetakToDB(status string) string {
	mapper := map[string]string{
		"NASABAH": "nasabah",
		"VENDOR":  "vendor",
		"BANK":    "bank",
	}
	if v, ok := mapper[status]; ok {
		return v
	}
	return "nasabah"
}
