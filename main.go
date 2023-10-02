package main

import (
	"database/sql"
	"fmt"

	"codemead.com/go_fintech/fintech_backend/api"
	"codemead.com/go_fintech/fintech_backend/api/handlers"
	"codemead.com/go_fintech/fintech_backend/api/repositories"
	db "codemead.com/go_fintech/fintech_backend/db/sqlc"
	"codemead.com/go_fintech/fintech_backend/token"
	"codemead.com/go_fintech/fintech_backend/utils"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		panic(fmt.Sprintf("Could not load env config: %v", err))
	}

	conn, err := sql.Open(config.DBdriver, config.DB_Source_Live)
	if err != nil {
		panic(fmt.Sprintf("Could not connect to database: %v", err))
	}
	queries := db.NewStore(conn)

	// fiber setup
	validate := validator.New()
	app := fiber.New()
	app.Use(cors.New())

	// validators
	err = validate.RegisterValidation("currency", currencyValidator)
	if err != nil {
		panic(fmt.Sprintf("connot register currency validator %s", err))
	}

	// token maker
	tokenMaker, err := token.NewPasetoMaker(config.SymetricKey)
	if err != nil {
		panic(fmt.Sprintf("cannot create token maker %w", err))
	}

	// init repositories
	authRepository := repositories.NewAuthRepository(queries, config, tokenMaker)
	userRepository := repositories.NewUserRepository(queries, config, tokenMaker)
	accountRepository := repositories.NewAccountRepository(queries, config, tokenMaker)

	// init handlers
	authHandlers := handlers.NewAuthHandlers(authRepository, validate)
	userHandlers := handlers.NewUserHandlers(userRepository, validate)
	accountHandlers := handlers.NewAccountHandlers(accountRepository, validate)

	server := api.NewServer(
		authHandlers,
		userHandlers,
		accountHandlers,
	)

	server.Fiber = app
	server.Config = config
	server.TokenMaker = tokenMaker
	server.Start()
}

var currencyValidator validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return utils.IsValidCurrency(currency)
	}
	return false
}
