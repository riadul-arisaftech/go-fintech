package handlers

import (
	"net/http"

	"codemead.com/go_fintech/fintech_backend/api/interfaces"
	"codemead.com/go_fintech/fintech_backend/api/models"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

type AuthHandlers struct {
	authRepository interfaces.IAuthRepository
	validator      *validator.Validate
}

func NewAuthHandlers(authRepository interfaces.IAuthRepository, validate *validator.Validate) *AuthHandlers {
	return &AuthHandlers{
		authRepository: authRepository,
		validator:      validate,
	}
}

func (h *AuthHandlers) Login(ctx *fiber.Ctx) error {
	var request models.RequestUserParams
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	if err := h.validator.Struct(request); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	authUser, err := h.authRepository.Login(request)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return ctx.Status(http.StatusOK).JSON(authUser)
}

func (h *AuthHandlers) Register(ctx *fiber.Ctx) error {
	var request models.RequestUserParams
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.validator.Struct(request); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	user, err := h.authRepository.Register(request)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "user already exist"})
			}
		}
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusCreated).JSON(user)
}
