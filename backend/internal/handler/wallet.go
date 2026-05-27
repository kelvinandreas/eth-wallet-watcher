package handler

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/constant"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/helper"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/model/request"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/model/response"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/service"
	"gorm.io/gorm"
)

type WalletHandler struct {
	walletService *service.WalletService
}

func NewWalletHandler(walletService *service.WalletService) *WalletHandler {
	return &WalletHandler{walletService: walletService}
}

func (h *WalletHandler) AddWallet(c *fiber.Ctx) error {
	var req request.WalletRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, constant.ErrInvalidRequestBody)
	}

	req.Address = strings.TrimSpace(req.Address)
	req.Label = strings.TrimSpace(req.Label)
	if req.Address == "" {
		return response.Error(c, fiber.StatusBadRequest, constant.ErrWalletAddressRequired)
	}
	if !helper.IsValidEthAddress(req.Address) {
		return response.Error(c, fiber.StatusBadRequest, constant.ErrInvalidEthAddress)
	}

	parsedUserID, err := helper.GetUserIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, err.Error())
	}

	if err := h.walletService.CreateWallet(parsedUserID, req.Address, req.Label); err != nil {
		if errors.Is(err, service.ErrWalletAlreadyExists) {
			return response.Error(c, fiber.StatusConflict, constant.ErrWalletAlreadyExists)
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusCreated, constant.MsgWalletCreated, nil)
}

func (h *WalletHandler) GetWallets(c *fiber.Ctx) error {
	parsedUserID, err := helper.GetUserIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, err.Error())
	}

	wallets, err := h.walletService.GetWalletsByUserID(parsedUserID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	walletResponses := make([]response.WalletResponse, 0, len(wallets))
	for _, wallet := range wallets {
		walletResponses = append(walletResponses, response.NewWalletResponse(wallet))
	}

	return response.Success(c, fiber.StatusOK, constant.MsgWalletRetrieved, walletResponses)
}

func (h *WalletHandler) DeleteWallet(c *fiber.Ctx) error {
	walletIDParam := strings.TrimSpace(c.Params("walletID"))
	walletID, err := uuid.Parse(walletIDParam)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, constant.ErrInvalidWalletID)
	}

	parsedUserID, err := helper.GetUserIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, err.Error())
	}

	if err := h.walletService.DeleteWallet(parsedUserID, walletID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(c, fiber.StatusNotFound, constant.ErrWalletNotFound)
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusOK, constant.MsgWalletDeleted, nil)
}
