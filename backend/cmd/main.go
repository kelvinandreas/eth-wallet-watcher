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
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	config.Init()

	if err := infrastructure.InitDB(); err != nil {
		log.Fatal("DB error:", err)
	}

	// Wire dependencies
	userRepo := repository.NewUserRepository(infrastructure.DB)
	walletRepo := repository.NewWalletRepository(infrastructure.DB)

	authService := service.NewAuthService(userRepo)
	walletService := service.NewWalletService(walletRepo)

	authHandler := handler.NewAuthHandler(authService)
	walletHandler := handler.NewWalletHandler(walletService)

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

	log.Fatal(app.Listen(":8080"))
}
