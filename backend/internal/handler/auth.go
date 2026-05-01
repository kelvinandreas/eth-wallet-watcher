package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/constant"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/model/request"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/model/response"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req request.AuthRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, constant.ErrInvalidRequestBody)
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
	if req.Email == "" || req.Password == "" {
		return response.Error(c, fiber.StatusBadRequest, constant.ErrEmailAndPasswordRequired)
	}

	if err := h.authService.Register(req.Email, req.Password); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusCreated, constant.MsgRegistrationSuccessful, nil)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req request.AuthRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, constant.ErrInvalidRequestBody)
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
	if req.Email == "" || req.Password == "" {
		return response.Error(c, fiber.StatusBadRequest, constant.ErrEmailAndPasswordRequired)
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, err.Error())
	}

	return response.Success(c, fiber.StatusOK, constant.MsgLoginSuccessful, response.NewAuthResponse(token))
}
