package interfaces

import (
	"codemead.com/go_fintech/fintech_backend/api/models"
	db "codemead.com/go_fintech/fintech_backend/db/sqlc"
	"github.com/gofiber/fiber/v2"
)

type IAccountRepository interface {
	GetUserAccounts(userId int64) ([]db.Account, error)
	CreateAccount(userId int64, request models.AccountRequest) (*db.Account, error)
	CreateTransfer(userId int64, request models.TransferRequest) (*db.TransferTxResponse, int, error)
}

type IAccountHandlers interface {
	GetUserAccounts(ctx *fiber.Ctx) error
	CreateAccount(ctx *fiber.Ctx) error
	CreateTransfer(ctx *fiber.Ctx) error
}
