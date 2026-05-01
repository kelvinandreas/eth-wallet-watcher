package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/config"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/handler"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/infrastructure"
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
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	app := fiber.New()

	// Auth routes (public)
	auth := app.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	log.Fatal(app.Listen(":8080"))
}
