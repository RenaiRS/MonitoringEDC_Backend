package dto

import "time"

type OverdueMesinResponse struct {
	TerminalID         string
	NamaNasabah        string
	LokasiMesin        string
	TanggalPasang      string `json:"tanggal_pasang"` // "2025-01-10"
	EstimasiPerbaikan  *time.Time
	TerlambatHari      int
	StatusPerbaikan    string
	Kerugian           int
}
