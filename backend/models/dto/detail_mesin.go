package dto

import "time"

type DetailMesinResponse struct {
	InformasiMesin struct {
		TerminalID  string
		MID         string
		TipeEDC     string
		Vendor      string
		StatusMesin string
		StatusData  string
	} `json:"informasi_mesin"`

	InformasiLokasi struct {
		NamaNasabah    string
		Cabang         string
		Kota           string
		TanggalPasang  time.Time
	} `json:"informasi_lokasi"`

	InformasiSewa struct {
		StatusSewa   string
		BiayaBulanan int
		EstimasiSelesai *time.Time
	} `json:"informasi_sewa"`

	SumberData string `json:"sumber_data"`
}
