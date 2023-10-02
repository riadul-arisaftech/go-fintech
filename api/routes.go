package api

import (
	"codemead.com/go_fintech/fintech_backend/api/middlewares"
)

func (s *Server) Routes() {
	v1 := s.Fiber.Group("/v1")

	// auth routes
	authRoutes := v1.Group("/auth")
	authRoutes.Post("/login", s.authHandlers.Login)
	authRoutes.Post("/register", s.authHandlers.Register)

	// user routes
	userRoutes := v1.Group("/users")
	userRoutes.Use(middlewares.AuthMiddleware(s.TokenMaker))
	userRoutes.Get("/", s.userHandlers.ListUsers)
	userRoutes.Get("/profile", s.userHandlers.GetProfile)

	// account routes
	accountRoutes := v1.Group("/account")
	accountRoutes.Use(middlewares.AuthMiddleware(s.TokenMaker))
	accountRoutes.Get("/", s.accHandlers.GetUserAccounts)
	accountRoutes.Post("/create", s.accHandlers.CreateAccount)
	accountRoutes.Post("/transfer", s.accHandlers.CreateTransfer)
}
