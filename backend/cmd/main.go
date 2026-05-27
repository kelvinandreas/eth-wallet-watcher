package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/config"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/handler"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/infrastructure"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/middleware"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/repository"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/service"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/worker"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	config.Init()

	if err := infrastructure.InitDB(); err != nil {
		log.Fatal("DB error:", err)
	}
	defer infrastructure.CloseDB()

	if err := infrastructure.InitRedis(); err != nil {
		log.Fatal("Redis error:", err)
	}
	defer infrastructure.CloseRedis()

	// Wire dependencies
	userRepo := repository.NewUserRepository(infrastructure.DB)
	walletRepo := repository.NewWalletRepository(infrastructure.DB)
	txRepo := repository.NewTransactionRepository(infrastructure.DB)
	notifRepo := repository.NewNotificationRepository(infrastructure.DB)

	etherscanSvc := service.NewEtherscanService(config.AppConfig.EtherscanKey, config.AppConfig.EtherscanBaseURL)

	authService := service.NewAuthService(userRepo)
	walletService := service.NewWalletService(walletRepo)
	txService := service.NewTransactionService(txRepo, infrastructure.RedisClient)
	notifService := service.NewNotificationService(notifRepo, infrastructure.RedisClient)

	authHandler := handler.NewAuthHandler(authService)
	walletHandler := handler.NewWalletHandler(walletService)
	txHandler := handler.NewTransactionHandler(txService)
	notifHandler := handler.NewNotificationHandler(notifService)

	// Background worker
	pollHandler := worker.NewPollWalletsHandler(walletRepo, txRepo, notifRepo, etherscanSvc, infrastructure.RedisClient)
	asynqServer := worker.NewServer()
	go worker.StartServer(asynqServer, pollHandler)
	go worker.StartScheduler()

	app := fiber.New()

	// Auth routes (public)
	auth := app.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Wallet routes (protected)
	wallet := app.Group("/wallets", middleware.JWTProtected())
	wallet.Post("", walletHandler.AddWallet)
	wallet.Get("", walletHandler.GetWallets)
	wallet.Delete(":walletID", walletHandler.DeleteWallet)
	wallet.Get(":walletID/transactions", txHandler.GetByWallet)

	// Notification routes (protected)
	notif := app.Group("/notifications", middleware.JWTProtected())
	notif.Get("", notifHandler.GetNotifications)
	notif.Patch(":notifID/read", notifHandler.MarkAsRead)

	log.Fatal(app.Listen(":8080"))
}
