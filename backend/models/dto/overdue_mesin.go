package dto

import "time"

type OverdueMesinResponse struct {
	TerminalID         string
	NamaNasabah        string
	LokasiMesin        string
	TanggalPasang      time.Time
	EstimasiPerbaikan  *time.Time
	TerlambatHari      int
	StatusPerbaikan    string
	Kerugian           int
}
