package handlers

import (
	"net/http"

	"codemead.com/go_fintech/fintech_backend/api/interfaces"
	"codemead.com/go_fintech/fintech_backend/api/models"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

type AccountHandlers struct {
	accountRepository interfaces.IAccountRepository
	validator         *validator.Validate
}

func NewAccountHandlers(accRepository interfaces.IAccountRepository, validate *validator.Validate) *AccountHandlers {
	return &AccountHandlers{
		accountRepository: accRepository,
		validator:         validate,
	}
}

func (h *AccountHandlers) GetUserAccounts(ctx *fiber.Ctx) error {
	userId := ctx.Locals("user_id").(int64)
	if userId == 0 {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "No authorized to access resources"})
	}

	accounts, err := h.accountRepository.GetUserAccounts(userId)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return ctx.Status(http.StatusOK).JSON(accounts)
}

func (h *AccountHandlers) CreateAccount(ctx *fiber.Ctx) error {
	userId := ctx.Locals("user_id").(int64)
	if userId == 0 {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "No authorized to access resources"})
	}

	var acc models.AccountRequest
	if err := ctx.BodyParser(&acc); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	if err := h.validator.Struct(acc); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	account, err := h.accountRepository.CreateAccount(userId, acc)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "already have and account with this currency"})
			}
		}
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return ctx.Status(http.StatusCreated).JSON(account)
}

func (h *AccountHandlers) CreateTransfer(ctx *fiber.Ctx) error {
	userId := ctx.Locals("user_id").(int64)
	if userId == 0 {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "No authorized to access resources"})
	}

	var tr models.TransferRequest
	if err := ctx.BodyParser(&tr); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	transfer, status, err := h.accountRepository.CreateTransfer(userId, tr)
	if err != nil {
		return ctx.Status(status).JSON(fiber.Map{"message": err.Error()})
	}

	return ctx.Status(status).JSON(transfer)
}
