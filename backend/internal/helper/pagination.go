package helper

import "github.com/gofiber/fiber/v2"

const (
	defaultPage  = 1
	defaultLimit = 10
	maxLimit     = 100
)

func GetPaginationParams(c *fiber.Ctx) (page, limit int) {
	page = c.QueryInt("page", defaultPage)
	limit = c.QueryInt("limit", defaultLimit)

	if page < 1 {
		page = defaultPage
	}
	if limit < 1 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	return page, limit
}

func GetOffset(page, limit int) int {
	return (page - 1) * limit
}
