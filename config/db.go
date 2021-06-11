package config

import (
	"github/sing3demons/covid-self-monitoring/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitializeDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("db.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// db.Migrator().DropTable(&models.Measurement{})
	db.AutoMigrate(&models.Measurement{}, &models.Symptom{}, &models.User{})

	return db
}
