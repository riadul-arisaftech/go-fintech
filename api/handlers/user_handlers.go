package handlers

import (
	"net/http"

	"codemead.com/go_fintech/fintech_backend/api/interfaces"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type UserHandlers struct {
	userRepository interfaces.IUserRepository
	validator      *validator.Validate
}

func NewUserHandlers(userRepository interfaces.IUserRepository, validate *validator.Validate) *UserHandlers {
	return &UserHandlers{
		userRepository: userRepository,
		validator:      validate,
	}
}

func (h *UserHandlers) ListUsers(ctx *fiber.Ctx) error {
	users, err := h.userRepository.ListUsers()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(http.StatusOK).JSON(users)
}

func (h *UserHandlers) GetProfile(ctx *fiber.Ctx) error {
	userId := ctx.Locals("user_id").(int64)
	if userId == 0 {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "No authorized to access resources"})
	}

	user, err := h.userRepository.GetProfile(userId)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return ctx.Status(http.StatusOK).JSON(user)
}
