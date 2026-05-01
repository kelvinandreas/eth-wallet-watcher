package helper

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/constant"
)

func GetUserIDFromLocals(c *fiber.Ctx) (uuid.UUID, error) {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, constant.ErrInvalidUserID)
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, constant.ErrInvalidUserID)
	}

	return parsedUserID, nil
}

func BuildURLWithQuery(baseURL string, params map[string]string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	q := u.Query()
	for key, value := range params {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}
