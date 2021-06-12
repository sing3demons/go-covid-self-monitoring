package controller

import (
	"github/sing3demons/covid-self-monitoring/models"
	"github/sing3demons/covid-self-monitoring/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type MeasurementController interface {
	Find(ctx *fiber.Ctx) error
	Create(ctx *fiber.Ctx) error
}

type measurementController struct {
	DB *gorm.DB
}

type creatMeasurement struct {
	Temperature float64 `json:"temperature" form:"temperature" validate:"required"`
	O2sat       int     `json:"o2sat" form:"o2sat" validate:"required"`
	Systolic    int     `json:"systolic" form:"systolic" validate:"required"`
	Diastolic   int     `json:"diastolic" form:"diastolic" validate:"required"`
	SymptomID   []uint  `json:"symptom_id" form:"symptom_id" validate:"required"`
}

type measurementResponse struct {
	CreatedAt   time.Time `json:"created_at"`
	Temperature float64   `json:"temperature" `
	O2sat       int       `json:"o2sat" `
	Systolic    int       `json:"systolic"`
	Diastolic   int       `json:"diastolic"`

	UserID uint `json:"user_id"`
	User   struct {
		Name string `json:"name"`
	} `json:"user"`

	Symptom []struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"symptom"`
}

func NewMeasurementController(db *gorm.DB) MeasurementController {
	return &measurementController{DB: db}
}

func (m *measurementController) Find(ctx *fiber.Ctx) error {
	measurement := []models.Measurement{}
	if err := m.DB.Preload("User").Preload("Symptom").Order("id desc").Find(&measurement).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err})
	}

	serializedMeasurement := []measurementResponse{}
	copier.Copy(&serializedMeasurement, &measurement)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": serializedMeasurement})
}

func (m *measurementController) Create(ctx *fiber.Ctx) error {
	var form creatMeasurement
	if err := ctx.BodyParser(&form); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	user, _ := utils.GetLocalsUser(ctx)

	symptomID := form.SymptomID
	symptoms := make([]models.Symptom, len(symptomID))
	for index, id := range symptomID {

		symptoms[index] = models.Symptom{ID: id}
	}

	var measurement = models.Measurement{Symptom: symptoms, UserID: user.ID}
	copier.Copy(&measurement, &form)

	if err := m.DB.Create(&measurement).Error; err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	var serializedMeasurement measurementResponse
	copier.Copy(&serializedMeasurement, &measurement)
	return ctx.Status(fiber.StatusCreated).JSON(serializedMeasurement)
}
