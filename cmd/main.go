package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/database"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := database.InitDB(); err != nil {
		log.Fatal("DB error:", err)
	}

	app := fiber.New()

	log.Fatal(app.Listen(":8080"))
}
