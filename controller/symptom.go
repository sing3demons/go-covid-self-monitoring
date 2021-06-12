package controller

import (
	"github/sing3demons/covid-self-monitoring/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type SymptomController interface {
	Create(ctx *fiber.Ctx) error
	FindAll(ctx *fiber.Ctx) error
	FindOne(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
}

type symptomController struct {
	DB *gorm.DB
}

type symptomResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type symptomForm struct {
	Name string `json:"name" form:"name" validate:"required"`
}

func NewSymptomController(DB *gorm.DB) SymptomController {
	return &symptomController{DB: DB}
}

func (s *symptomController) FindAll(ctx *fiber.Ctx) error {
	symptom := []models.Symptom{}

	s.DB.Find(&symptom)

	serializedSymptom := []symptomResponse{}
	copier.Copy(&serializedSymptom, &symptom)
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"symptom": serializedSymptom})
}

func (s *symptomController) FindOne(ctx *fiber.Ctx) error {
	symptom, err := s.findSymptomByID(ctx)
	if err != nil {
		ctx.Status(fiber.StatusNotFound)
		return ctx.JSON(fiber.Map{"error": err.Error()})
	}

	var serializedSymptom symptomResponse
	copier.Copy(&serializedSymptom, &symptom)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"symptom": serializedSymptom})

}

func (s *symptomController) Create(ctx *fiber.Ctx) error {
	var form symptomForm
	if err := ctx.BodyParser(&form); err != nil {
		if err := ctx.BodyParser(&form); err != nil {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}
	}

	var symptom models.Symptom
	copier.Copy(&symptom, &form)
	if err := s.DB.Create(&symptom).Error; err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	var serializedSymptom symptomResponse
	copier.Copy(&serializedSymptom, &symptom)

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"symptom": serializedSymptom})
}

//Update - put method
func (s *symptomController) Update(ctx *fiber.Ctx) error {
	symptom, err := s.findSymptomByID(ctx)
	if err != nil {
		ctx.Status(fiber.StatusNotFound)
		return ctx.JSON(fiber.Map{"error": err.Error()})
	}

	var form symptomForm
	if err := ctx.BodyParser(&form); err != nil {
		ctx.Status(fiber.StatusUnprocessableEntity)
		return ctx.JSON(fiber.Map{"error": err.Error()})
	}

	copier.Copy(&symptom, &form)
	if err := s.DB.Model(&symptom).Updates(&symptom).Error; err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (s *symptomController) Delete(ctx *fiber.Ctx) error {
	symptom, err := s.findSymptomByID(ctx)
	if err != nil {
		ctx.Status(fiber.StatusNotFound)
		return ctx.JSON(fiber.Map{"error": err.Error()})
	}

	if err := s.DB.Delete(&symptom).Error; err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(fiber.StatusNoContent)

}

func (s *symptomController) findSymptomByID(ctx *fiber.Ctx) (*models.Symptom, error) {
	var symptom models.Symptom
	id := ctx.Params("id")

	if err := s.DB.First(&symptom, id).Error; err != nil {
		return nil, err
	}

	return &symptom, nil
}
