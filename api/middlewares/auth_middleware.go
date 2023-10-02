package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"codemead.com/go_fintech/fintech_backend/token"
	"github.com/gofiber/fiber/v2"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
)

func AuthMiddleware(tokenMaker token.Maker) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authorizationHeader := ctx.Get(AuthorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != AuthorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s == %s", authorizationType, AuthorizationTypeBearer)
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
		}

		accessToken := fields[1]
		userId, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
		}

		ctx.Locals("user_id", userId)
		return ctx.Next()
	}
}
