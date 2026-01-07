package dto

import "time"

type RekapMesinResponse struct {
	ID           uint      `json:"id"`
	TerminalID   string    `json:"terminal_id"`
	NamaNasabah  string    `json:"nama_nasabah"`
	Cabang       string    `json:"cabang"`
	TanggalPasang time.Time `json:"tanggal_pasang"`
	StatusMesin  string    `json:"status_mesin"`
	StatusData   string    `json:"status_data"`
}
