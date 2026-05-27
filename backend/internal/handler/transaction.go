package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/constant"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/helper"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/model/response"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/service"
)

type TransactionHandler struct {
	txService *service.TransactionService
}

func NewTransactionHandler(txService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{txService: txService}
}

func (h *TransactionHandler) GetByWallet(c *fiber.Ctx) error {
	walletID, err := uuid.Parse(strings.TrimSpace(c.Params("walletID")))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, constant.ErrInvalidWalletID)
	}

	page, limit := helper.GetPaginationParams(c)

	txs, total, err := h.txService.GetByWalletID(walletID, page, limit)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	items := make([]response.TransactionResponse, 0, len(txs))
	for _, tx := range txs {
		items = append(items, response.NewTransactionResponse(tx))
	}

	return response.Success(c, fiber.StatusOK, constant.MsgTransactionsRetrieved, response.NewPaginatedResponse(items, total, page, limit))
}
