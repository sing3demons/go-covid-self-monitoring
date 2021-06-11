package controller

import (
	"github/sing3demons/covid-self-monitoring/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type SymptomController interface {
	Create(ctx *fiber.Ctx) error
	Find(ctx *fiber.Ctx) error
}

type symptomController struct {
	DB *gorm.DB
}

type symptomResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type createSymptom struct {
	Name string `json:"name" form:"name" validate:"required"`
}

func NewSymptomController(DB *gorm.DB) SymptomController {
	return &symptomController{DB: DB}
}

func (s *symptomController) Find(ctx *fiber.Ctx) error {
	symptom := []models.Symptom{}

	s.DB.Find(&symptom)

	serializedSymptom := []symptomResponse{}
	copier.Copy(&serializedSymptom, &symptom)
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"symptom": serializedSymptom})
}

func (s *symptomController) Create(ctx *fiber.Ctx) error {
	var form createSymptom
	if err := ctx.BodyParser(&form); err != nil {
		if err := ctx.BodyParser(&form); err != nil {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}
	}

	var symptom models.Symptom
	copier.Copy(&symptom, &form)
	s.DB.Create(&symptom)

	var serializedSymptom symptomResponse
	copier.Copy(&serializedSymptom, &symptom)

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"symptom": serializedSymptom})
}
