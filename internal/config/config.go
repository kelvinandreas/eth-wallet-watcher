package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	DBHost       string
	DBUser       string
	DBPassword   string
	DBName       string
	DBPort       string
	RedisAddr    string
	EtherscanKey string
	ResendKey    string
	JWTSecret    string
	TokenTTL     time.Duration
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

var AppConfig *Config

func (c *Config) GetDBDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort,
	)
}

func Init() {
	hours, err := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	if err != nil || hours == 0 {
		hours = 24
	}

	AppConfig = &Config{
		DBHost:       os.Getenv("DB_HOST"),
		DBUser:       os.Getenv("DB_USER"),
		DBPassword:   os.Getenv("DB_PASSWORD"),
		DBName:       os.Getenv("DB_NAME"),
		DBPort:       os.Getenv("DB_PORT"),
		RedisAddr:    os.Getenv("REDIS_ADDR"),
		EtherscanKey: os.Getenv("ETHERSCAN_API_KEY"),
		ResendKey:    os.Getenv("RESEND_API_KEY"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
		TokenTTL:     time.Duration(hours) * time.Hour,
	}
}
