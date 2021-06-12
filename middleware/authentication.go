package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

func Authenticate() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:     []byte(os.Getenv("JWT_SECRET_KEY")),
		ErrorHandler:   jwtError,
		ContextKey:     "user",
		SuccessHandler: MiddlewareFunc(),
	})
}

func MiddlewareFunc() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}
