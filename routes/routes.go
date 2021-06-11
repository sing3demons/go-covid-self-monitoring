package routes

import (
	"github/sing3demons/covid-self-monitoring/config"
	"github/sing3demons/covid-self-monitoring/controller"

	"github.com/gofiber/fiber/v2"
)

func Serve(app *fiber.App) {
	db := config.InitializeDB()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON("hello, world")
	})

	v1 := app.Group("/api/v1/")

	measurementController := controller.NewMeasurementController(db)
	measurementGroup := v1.Group("measurement")
	{
		measurementGroup.Get("", measurementController.Find)
		measurementGroup.Post("", measurementController.Create)
	}

	symptomController := controller.NewSymptomController(db)
	symptomGroup := v1.Group("symptom")
	{
		symptomGroup.Get("", symptomController.Find)
		symptomGroup.Post("", symptomController.Create)
	}
}