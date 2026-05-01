package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/config"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/constant"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/model/response"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Error(c, fiber.StatusUnauthorized, constant.ErrMissingAuthorizationHeader)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.Error(c, fiber.StatusUnauthorized, constant.ErrInvalidAuthorizationFormat)
		}

		tokenString := parts[1]

		token, err := jwt.ParseWithClaims(tokenString, &config.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return response.Error(c, fiber.StatusUnauthorized, constant.ErrInvalidOrExpiredToken)
		}

		claims, ok := token.Claims.(*config.Claims)
		if !ok {
			return response.Error(c, fiber.StatusUnauthorized, constant.ErrInvalidTokenClaims)
		}

		c.Locals("userID", claims.UserID)

		return c.Next()
	}
}
