// models/overdue.go
package models

import "time"

type Overdue struct {
	ID                 uint      `gorm:"primaryKey"`
	MesinID            uint
	Mesin              MesinEDC  `gorm:"foreignKey:MesinID"`

	StatusPerbaikan    string    // perbaikan | warning | overdue
	EstimasiSelesai    time.Time
	TerlambatHari      int
	EstimasiKerugian   int

	CreatedAt          time.Time
	UpdatedAt          time.Time
}
