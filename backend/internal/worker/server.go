package worker

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/config"
)

func NewServer() *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{Addr: config.AppConfig.RedisAddr},
		asynq.Config{
			Concurrency: 5,
			Queues:      map[string]int{"default": 1},
		},
	)
}

func StartServer(srv *asynq.Server, handler *PollWalletsHandler) {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskPollWallets, handler.ProcessTask)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("asynq server error: %v", err)
	}
}

func StartScheduler() {
	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{Addr: config.AppConfig.RedisAddr},
		nil,
	)

	if _, err := scheduler.Register("*/5 * * * *", asynq.NewTask(TaskPollWallets, nil)); err != nil {
		log.Fatalf("register scheduler: %v", err)
	}

	if err := scheduler.Run(); err != nil {
		log.Fatalf("asynq scheduler error: %v", err)
	}
}
