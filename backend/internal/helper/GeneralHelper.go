package helper

import (
	"math/big"
	"net/url"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/constant"
)

var weiPerEth = new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))

func WeiToETH(weiStr string) string {
	wei := new(big.Int)
	if _, ok := wei.SetString(weiStr, 10); !ok {
		return "0"
	}
	eth := new(big.Float).Quo(new(big.Float).SetInt(wei), weiPerEth)
	return eth.Text('f', 8)
}

var ethAddressRegex = regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)

func IsValidEthAddress(address string) bool {
	return ethAddressRegex.MatchString(address)
}

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
