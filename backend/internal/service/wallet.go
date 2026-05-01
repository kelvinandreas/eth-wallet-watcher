package service

import (
	"github.com/google/uuid"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/model/app"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/repository"
)

type WalletService struct {
	walletRepo *repository.WalletRepository
}

func NewWalletService(walletRepo *repository.WalletRepository) *WalletService {
	return &WalletService{walletRepo: walletRepo}
}

func (s *WalletService) CreateWallet(userID uuid.UUID, address, label string) error {
	wallet := &app.MonitoredWallet{
		UserID:  userID,
		Address: address,
		Label:   label,
	}

	return s.walletRepo.CreateWallet(wallet)
}

func (s *WalletService) GetWalletsByUserID(userID uuid.UUID) ([]app.MonitoredWallet, error) {
	return s.walletRepo.GetWalletsByUserID(userID)
}

func (s *WalletService) DeleteWallet(userID, walletID uuid.UUID) error {
	return s.walletRepo.DeleteWallet(userID, walletID)
}

func (s *WalletService) GetAllWallets() ([]app.MonitoredWallet, error) {
	return s.walletRepo.GetAllWallets()
}
