package database

import (
	"fmt"
	"log"
	"os"

	"backend/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBLoad() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect database:", err)
	}

	log.Println("✅ Database connected")
}

func DBMigrate() {
	if err := DB.AutoMigrate(
		&models.MesinEDC{},
		&models.Sewa{},
		&models.Perbaikan{},
		&models.Overdue{},
	); err != nil {
		log.Fatal("❌ Failed to migrate:", err)
	}
	log.Println("✅ Database migrated")
}


