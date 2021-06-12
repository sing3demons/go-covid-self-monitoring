package routes

import (
	"github/sing3demons/covid-self-monitoring/config"
	"github/sing3demons/covid-self-monitoring/controller"
	"github/sing3demons/covid-self-monitoring/middleware"

	"github.com/gofiber/fiber/v2"
)

func Serve(app *fiber.App) {
	db := config.GetDB()
	authenticate := middleware.Authenticate()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON("hello, world")
	})

	v1 := app.Group("/api/v1/")

	measurementController := controller.NewMeasurementController(db)
	measurementGroup := v1.Group("measurement")
	measurementGroup.Use(authenticate)
	{
		measurementGroup.Get("", measurementController.Find)
		measurementGroup.Post("", measurementController.Create)
	}

	symptomController := controller.NewSymptomController(db)
	symptomGroup := v1.Group("symptom")
	symptomGroup.Use(authenticate)
	{
		symptomGroup.Get("", symptomController.Find)
		symptomGroup.Post("", symptomController.Create)
	}

	authController := controller.NewAuthController(db)
	authGroup := v1.Group("auth")
	{
		authGroup.Post("/register", authController.Register)
		authGroup.Post("/login", authController.Login)
	}
	{
		authGroup.Use(authenticate)
		authGroup.Get("/profile", authController.Profile)
	}
}
