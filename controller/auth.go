package controller

import (
	"github/sing3demons/covid-self-monitoring/models"
	"github/sing3demons/covid-self-monitoring/utils"
	"os"
	"strconv"
	"time"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authForm struct {
	Name     string `json:"name" form:"name" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=6"`
}

type AuthController interface {
	Register(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	Profile(ctx *fiber.Ctx) error
}

type authController struct {
	DB *gorm.DB
}

type authResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func NewAuthController(db *gorm.DB) AuthController {
	return &authController{DB: db}
}

func (a *authController) Register(ctx *fiber.Ctx) error {
	var form authForm
	if err := ctx.BodyParser(&form); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	var user models.User
	copier.Copy(&user, &form)
	user.Password = user.GenerateEncryptedPassword()
	if err := a.DB.Create(&user).Error; err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	var serializedUser authResponse
	copier.Copy(&serializedUser, user)
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"user": serializedUser})
}

func (a *authController) Profile(ctx *fiber.Ctx) error {
	user, err := utils.GetLocalsUser(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	var serializedUser authResponse
	copier.Copy(&serializedUser, user)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"user": serializedUser})
}

func (a *authController) Login(ctx *fiber.Ctx) error {
	var form login
	if err := ctx.BodyParser(&form); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	var user models.User
	// copier.Copy(&user, &form)

	if err := a.DB.Where("email = ?", form.Email).First(&user).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	var serializedUser authResponse
	copier.Copy(&serializedUser, user)

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"token": token})
}
