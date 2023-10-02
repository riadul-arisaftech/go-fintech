package interfaces

import (
	"codemead.com/go_fintech/fintech_backend/api/models"
	"github.com/gofiber/fiber/v2"
)

type IUserRepository interface {
	ListUsers() ([]models.ResponseUser, error)
	GetProfile(userId int64) (*models.ResponseUser, error)
}

type IUserHandlers interface {
	ListUsers(ctx *fiber.Ctx) error
	GetProfile(ctx *fiber.Ctx) error
}
