package api

import (
	"fmt"

	"codemead.com/go_fintech/fintech_backend/api/interfaces"
	"codemead.com/go_fintech/fintech_backend/token"
	"codemead.com/go_fintech/fintech_backend/utils"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

type Server struct {
	Config     *utils.Config
	Fiber      *fiber.App
	TokenMaker token.Maker

	// all handlers interfaces
	authHandlers interfaces.IAuthHandlers
	userHandlers interfaces.IUserHandlers
	accHandlers  interfaces.IAccountHandlers
}

func NewServer(
	authHandlers interfaces.IAuthHandlers,
	uHandlers interfaces.IUserHandlers,
	accHandlers interfaces.IAccountHandlers,
) *Server {
	return &Server{
		authHandlers: authHandlers,
		userHandlers: uHandlers,
		accHandlers:  accHandlers,
	}
}

func (s *Server) Start() {
	s.Routes()
	s.Fiber.Listen(fmt.Sprintf(":%v", s.Config.ServerPort))
}
