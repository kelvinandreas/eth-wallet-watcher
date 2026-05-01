package infrastructure

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/config"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var AsynqClient *asynq.Client

func InitRedis() error {
	client := redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		DB:           0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return err
	}

	RedisClient = client
	AsynqClient = asynq.NewClient(asynq.RedisClientOpt{Addr: config.AppConfig.RedisAddr})
	return nil
}

func CloseRedis() {
	if AsynqClient != nil {
		AsynqClient.Close()
	}
	if RedisClient != nil {
		RedisClient.Close()
	}
}
