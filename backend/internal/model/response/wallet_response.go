package response

import (
	"github.com/google/uuid"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/model/app"
)

type WalletResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Address   string    `json:"address"`
	Label     string    `json:"label"`
	LastBlock string    `json:"last_block,omitempty"`
}

func NewWalletResponse(wallet app.MonitoredWallet) WalletResponse {
	return WalletResponse{
		ID:        wallet.ID,
		UserID:    wallet.UserID,
		Address:   wallet.Address,
		Label:     wallet.Label,
		LastBlock: wallet.LastBlock,
	}
}
