package dto

type DetailMesinResponse struct {
	InformasiMesin struct {
		TerminalID  string `json:"terminal_id"`
		MID         string `json:"mid"`
		TipeEDC     string `json:"tipe_edc"`
		Vendor      string `json:"vendor"`
		StatusMesin string `json:"status_mesin"`
		StatusData  string `json:"status_data"`
	} `json:"informasi_mesin"`

	InformasiLokasi struct {
		NamaNasabah   string `json:"nama_nasabah"`
		Cabang        string `json:"cabang"`
		Kota          string `json:"kota"`
		TanggalPasang string `json:"tanggal_pasang"` // Format: YYYY-MM-DD
		StatusLetak   string `json:"status_letak"`
	} `json:"informasi_lokasi"`

	InformasiSewa struct {
		StatusSewa      string  `json:"status_sewa"`
		BiayaBulanan    int     `json:"biaya_bulanan"`
		EstimasiSelesai *string `json:"estimasi_selesai,omitempty"` // Format: YYYY-MM-DD
	} `json:"informasi_sewa"`

	SumberData string `json:"sumber_data"`
}