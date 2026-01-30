package models

import "time"

type MesinEDC struct {
	ID           uint      `gorm:"primaryKey"`
	TerminalID   string    `gorm:"unique;not null"`
	MID          string

	NamaNasabah  string
	Kota         string
	Cabang       string

	TipeEDC      string
	Vendor       string

	StatusMesin  string    // aktif | perbaikan | rusak | tidak_aktif
	StatusData   string    // bank | vendor_only
	LetakMesin   string    // nasabah | bank | vendor

	TanggalPasang *time.Time `gorm:"column:tanggal_pasang"`

	CreatedAt    time.Time
	UpdatedAt    time.Time

	// RELATION (POINTER)
	Sewas []Sewa `gorm:"foreignKey:MesinID"`
	Perbaikan []Perbaikan  `gorm:"foreignKey:MesinID"`
}
