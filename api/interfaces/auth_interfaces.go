package interfaces

import (
	"codemead.com/go_fintech/fintech_backend/api/models"
	"github.com/gofiber/fiber/v2"
)

type IAuthRepository interface {
	Login(request models.RequestUserParams) (*models.ResponseAuthUser, error)
	Register(request models.RequestUserParams) (*models.ResponseUser, error)
}

type IAuthHandlers interface {
	Login(ctx *fiber.Ctx) error
	Register(ctx *fiber.Ctx) error
}
