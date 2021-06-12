package main

import (
	"github/sing3demons/covid-self-monitoring/config"
	"github/sing3demons/covid-self-monitoring/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.InitializeDB()

	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())

	routes.Serve(app)

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
