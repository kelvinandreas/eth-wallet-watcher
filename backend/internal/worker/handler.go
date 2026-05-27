package worker

import (
	"context"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/constant"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/helper"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/model/app"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/repository"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/service"
	"github.com/redis/go-redis/v9"
)

type PollWalletsHandler struct {
	walletRepo  *repository.WalletRepository
	txRepo      *repository.TransactionRepository
	notifRepo   *repository.NotificationRepository
	etherscan   *service.EtherscanService
	redisClient *redis.Client
}

func NewPollWalletsHandler(
	walletRepo *repository.WalletRepository,
	txRepo *repository.TransactionRepository,
	notifRepo *repository.NotificationRepository,
	etherscan *service.EtherscanService,
	redisClient *redis.Client,
) *PollWalletsHandler {
	return &PollWalletsHandler{
		walletRepo:  walletRepo,
		txRepo:      txRepo,
		notifRepo:   notifRepo,
		etherscan:   etherscan,
		redisClient: redisClient,
	}
}

func (h *PollWalletsHandler) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	wallets, err := h.walletRepo.GetAllWallets()
	if err != nil {
		return fmt.Errorf("get wallets: %w", err)
	}

	for _, wallet := range wallets {
		if err := h.pollWallet(ctx, wallet); err != nil {
			log.Printf("poll wallet %s: %v", wallet.Address, err)
		}
	}

	return nil
}

func (h *PollWalletsHandler) pollWallet(ctx context.Context, wallet app.MonitoredWallet) error {
	result, err := h.etherscan.FetchTransactions(wallet.Address, wallet.LastBlock)
	if err != nil {
		return err
	}

	if len(result.Transactions) == 0 {
		return nil
	}

	txs := make([]app.Transaction, 0, len(result.Transactions))
	for _, etx := range result.Transactions {
		txs = append(txs, app.Transaction{
			WalletID:    wallet.ID,
			TxHash:      etx.Hash,
			FromAddress: etx.From,
			ToAddress:   etx.To,
			Value:       etx.Value,
			Timestamp:   etx.TimeStamp,
		})
	}

	if err := h.txRepo.UpsertTransactions(txs); err != nil {
		return fmt.Errorf("upsert transactions: %w", err)
	}

	for _, etx := range result.Transactions {
		exists, err := h.notifRepo.ExistsByTxHash(wallet.ID, etx.Hash)
		if err != nil {
			log.Printf("check notification exists for tx %s: %v", etx.Hash, err)
			continue
		}
		if exists {
			continue
		}

		ethValue := helper.WeiToETH(etx.Value)
		notif := &app.Notification{
			UserID:   wallet.UserID,
			WalletID: wallet.ID,
			TxHash:   etx.Hash,
			Message:  fmt.Sprintf("New transaction on wallet %s: %s ETH", wallet.Address, ethValue),
		}
		if err := h.notifRepo.Create(notif); err != nil {
			log.Printf("create notification for tx %s: %v", etx.Hash, err)
		}
	}

	if err := h.walletRepo.UpdateLastBlock(wallet.ID, result.LastBlock); err != nil {
		return fmt.Errorf("update last block: %w", err)
	}

	helper.DeleteCache(ctx, h.redisClient,
		helper.CacheKey(constant.CacheKey.Transactions, wallet.ID),
		helper.CacheKey(constant.CacheKey.Notifications, wallet.UserID),
	)

	return nil
}
