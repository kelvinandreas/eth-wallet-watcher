package repository

import (
	"github.com/google/uuid"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/model/app"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) UpsertTransactions(txs []app.Transaction) error {
	if len(txs) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "tx_hash"}},
		DoNothing: true,
	}).Create(&txs).Error
}

func (r *TransactionRepository) GetByWalletID(walletID uuid.UUID) ([]app.Transaction, error) {
	var txs []app.Transaction
	result := r.db.Where("wallet_id = ?", walletID).Order("timestamp desc").Find(&txs)
	return txs, result.Error
}
