package models

import "time"

type Sewa struct {
	ID           uint      `gorm:"primaryKey"`

	MesinID      uint
	Mesin 		 *MesinEDC `gorm:"foreignKey:MesinID;references:ID"`


	StatusSewa   string    // aktif | belum_terdata | berakhir
	BiayaBulanan int

	CreatedAt    time.Time
	UpdatedAt    time.Time
}
