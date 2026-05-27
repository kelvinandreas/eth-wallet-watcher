package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/constant"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/helper"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/model/app"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/repository"
	"github.com/redis/go-redis/v9"
)

type TransactionService struct {
	txRepo      *repository.TransactionRepository
	redisClient *redis.Client
}

func NewTransactionService(txRepo *repository.TransactionRepository, redisClient *redis.Client) *TransactionService {
	return &TransactionService{txRepo: txRepo, redisClient: redisClient}
}

func (s *TransactionService) GetByWalletID(walletID uuid.UUID, page, limit int) ([]app.Transaction, int64, error) {
	key := helper.CacheKey(constant.CacheKey.Transactions, walletID)

	all, err := helper.GetOrSet(context.Background(), s.redisClient, key, helper.DefaultCacheTTL, func() ([]app.Transaction, error) {
		return s.txRepo.GetByWalletID(walletID)
	})
	if err != nil {
		return nil, 0, err
	}

	total := int64(len(all))
	start, end := paginate(len(all), page, limit)
	return all[start:end], total, nil
}

func paginate(total, page, limit int) (start, end int) {
	start = (page - 1) * limit
	if start >= total {
		return total, total
	}
	end = start + limit
	if end > total {
		end = total
	}
	return start, end
}
