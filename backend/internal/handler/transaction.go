package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/constant"
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

	txs, err := h.txService.GetByWalletID(walletID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	result := make([]response.TransactionResponse, 0, len(txs))
	for _, tx := range txs {
		result = append(result, response.NewTransactionResponse(tx))
	}

	return response.Success(c, fiber.StatusOK, constant.MsgTransactionsRetrieved, result)
}
