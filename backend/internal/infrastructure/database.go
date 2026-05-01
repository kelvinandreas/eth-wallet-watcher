package infrastructure

import (
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/config"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/model/app"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	db, err := gorm.Open(postgres.Open(config.AppConfig.GetDBDSN()), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db
	return DB.AutoMigrate(&app.User{}, &app.MonitoredWallet{}, &app.Notification{}, &app.Transaction{})
}

func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		return
	}
	sqlDB.Close()
}
