package utils

import (
	"fmt"
	"github/sing3demons/covid-self-monitoring/config"
	"github/sing3demons/covid-self-monitoring/models"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func GetLocalsUser(c *fiber.Ctx) (*models.User, error) {
	db := config.GetDB()
	var user models.User

	// id := ctx.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)["sub"]
	contextKey := c.Locals("user").(*jwt.Token)
	sub := contextKey.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", sub["sub"])

	if err := db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
