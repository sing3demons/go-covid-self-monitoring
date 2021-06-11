package models

import "gorm.io/gorm"

type Measurement struct {
	gorm.Model
	Temperature float64 `gorm:"not null"`
	O2sat       int     `gorm:"not null"`
	Systolic    int     `gorm:"not null"`
	Diastolic   int     `gorm:"not null"`
	// SymptomID   uint
	Symptom []Symptom `gorm:"many2many:measurement_symptom"`
}
