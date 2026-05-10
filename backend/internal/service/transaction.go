package service

import (
	"github.com/google/uuid"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/model/app"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/repository"
)

type TransactionService struct {
	txRepo *repository.TransactionRepository
}

func NewTransactionService(txRepo *repository.TransactionRepository) *TransactionService {
	return &TransactionService{txRepo: txRepo}
}

func (s *TransactionService) GetByWalletID(walletID uuid.UUID) ([]app.Transaction, error) {
	return s.txRepo.GetByWalletID(walletID)
}
