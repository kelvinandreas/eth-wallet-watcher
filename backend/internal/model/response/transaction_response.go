package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/model/app"
)

type TransactionResponse struct {
	ID          uuid.UUID `json:"id"`
	WalletID    uuid.UUID `json:"wallet_id"`
	TxHash      string    `json:"tx_hash"`
	FromAddress string    `json:"from_address"`
	ToAddress   string    `json:"to_address"`
	Value       string    `json:"value"`
	Timestamp   time.Time `json:"timestamp"`
}

func NewTransactionResponse(tx app.Transaction) TransactionResponse {
	return TransactionResponse{
		ID:          tx.ID,
		WalletID:    tx.WalletID,
		TxHash:      tx.TxHash,
		FromAddress: tx.FromAddress,
		ToAddress:   tx.ToAddress,
		Value:       tx.Value,
		Timestamp:   tx.Timestamp,
	}
}
