package repository

import (
	"github.com/google/uuid"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/model/app"
	"gorm.io/gorm"
)

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) CreateWallet(wallet *app.MonitoredWallet) error {
	return r.db.Create(wallet).Error
}

func (r *WalletRepository) GetWalletsByUserID(userID uuid.UUID) ([]app.MonitoredWallet, error) {
	var wallets []app.MonitoredWallet
	result := r.db.Where("user_id = ?", userID).Find(&wallets)
	if result.Error != nil {
		return nil, result.Error
	}

	return wallets, nil
}

func (r *WalletRepository) DeleteWallet(userID uuid.UUID, walletID uuid.UUID) error {
	result := r.db.Delete(&app.MonitoredWallet{}, "id = ? AND user_id = ?", walletID, userID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *WalletRepository) GetAllWallets() ([]app.MonitoredWallet, error) {
	var wallets []app.MonitoredWallet
	result := r.db.Find(&wallets)
	if result.Error != nil {
		return nil, result.Error
	}

	return wallets, nil
}
