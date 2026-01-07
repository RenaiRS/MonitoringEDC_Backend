package models

import "time"

type Perbaikan struct {
	ID uint `gorm:"primaryKey"`

	MesinID uint
	Mesin   *MesinEDC `gorm:"foreignKey:MesinID"`

	StatusPerbaikan   string // perbaikan | warning | overdue
	EstimasiPerbaikan *time.Time
	EstimasiSelesai   *time.Time
	SN                string

	CreatedAt time.Time
	UpdatedAt time.Time
}
