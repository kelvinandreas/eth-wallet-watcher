package response

import "github.com/gofiber/fiber/v2"

type BaseResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
}

func Success(c *fiber.Ctx, statusCode int, message string, data any) error {
	return c.Status(statusCode).JSON(BaseResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	})
}

func Error(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(BaseResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       nil,
	})
}
