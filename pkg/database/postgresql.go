package database

import (
	"errors"
	"fmt"
	"time"
	"w2learn/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
)

type PostgresConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string
	SSLMode      string
	MaxIdleConns int
	MaxOpenConns int
	MaxLifeTime  int
	LogLevel     string
}

var globalDB *gorm.DB

func NewPostgresDB(cfg *PostgresConfig) (*gorm.DB, error) {
	logger.Info("Check config Start")
	if cfg == nil {
		return nil, errors.New("config is nil")
	}
	logger.Info("Check config End")

	logger.Info("Prepare DSN Start")
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode, "Asia/Shanghai",
	)
	logger.Info("Prepare DSN End: ", zap.String("dsn", dsn))

	logger.Info("Check LogLevel Start")
	var logLevel logger2.LogLevel

	switch cfg.LogLevel {
	case "silent":
		logLevel = logger2.Silent
	case "error":
		logLevel = logger2.Error
	case "warn":
		logLevel = logger2.Warn
	case "info":
		logLevel = logger2.Info
	default:
		logLevel = logger2.Info
	}
	logger.Info("Check LogLevel Start")

	logger.Info("Connect to PostgreSQL Start")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:  logger2.Default.LogMode(logLevel),
		NowFunc: func() time.Time { return time.Now() },
	})

	if err != nil {
		logger.Error("Connect to PostgreSQL Error", zap.Error(err))
		return nil, err
	}

	baseDB, err := db.DB()

	if err != nil {
		logger.Error("Get PostgreSQL DB Error", zap.Error(err))
	}

	baseDB.SetMaxIdleConns(cfg.MaxIdleConns)
	baseDB.SetMaxOpenConns(cfg.MaxOpenConns)
	if cfg.MaxLifeTime > 0 {
		baseDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifeTime) * time.Second)
	}

	err = baseDB.Ping()

	if err != nil {
		logger.Error("Ping to PostgreSQL Error", zap.Error(err))
	}

	logger.Info("Connect to PostgreSQL End")

	globalDB = db

	return globalDB, nil
}

func GetDB() *gorm.DB {
	return globalDB
}

func Close() error {
	if globalDB != nil {
		baseDB, err := globalDB.DB()

		if err != nil {
			return err
		}

		return baseDB.Close()
	}

	return nil
}
