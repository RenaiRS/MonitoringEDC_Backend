package dto

type RekapMesinResponse struct {
	ID            uint   `json:"id"`
	TerminalID    string `json:"terminal_id"`
	NamaNasabah   string `json:"nama_nasabah"`
	Cabang        string `json:"cabang"`
	TanggalPasang string `json:"tanggal_pasang"` // YYYY-MM-DD
	StatusMesin   string `json:"status_mesin"`
	StatusData    string `json:"status_data"`
}
