package models

import "time"

type Perbaikan struct {
	ID uint `gorm:"primaryKey"`

	MesinID uint
	Mesin   *MesinEDC `gorm:"foreignKey:MesinID;references:ID"`

	StatusPerbaikan   string // perbaikan | warning | overdue
	EstimasiSelesai   *time.Time
	SN                string

	CreatedAt time.Time
	UpdatedAt time.Time
}
