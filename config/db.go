package config

import (
	"github/sing3demons/covid-self-monitoring/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitializeDB() {
	database, err := gorm.Open(sqlite.Open("db.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// db.Migrator().DropTable(&models.Measurement{}, &models.User{})
	database.AutoMigrate(&models.Measurement{}, &models.Symptom{}, &models.User{})

	db = database
}
func GetDB() *gorm.DB {
	return db
}
