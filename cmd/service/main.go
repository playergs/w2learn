package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"w2learn/internal/config"
	"w2learn/internal/controller"
	"w2learn/internal/model"
	"w2learn/internal/repository"
	"w2learn/internal/router"
	"w2learn/internal/service"
	"w2learn/internal/utils"
	"w2learn/pkg/database"
	"w2learn/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	log.Println("Load Service Type Start")
	serviceType := os.Getenv(config.ServiceTypeLabel)

	switch serviceType {
	case config.ServiceTypeProd:
		log.Println("Load Service Type is PROD")
	case config.ServiceTypeDev:
		log.Println("Load Service Type is DEV")
	default:
		serviceType = config.ServiceTypeDev
		log.Println("Load Service Type is DEV (DEFAULT)")
	}
	log.Println("Load Service Type End")

	log.Println("Load ConfigDir Start")
	configDir := os.Getenv(config.FileDirLabel)

	if configDir == "" {
		configDir = config.FileDirDefault
	}

	log.Println("Load ConfigDir End")

	log.Println("Load ConfigType Start")
	configFilePostfix := os.Getenv(config.FilePostfixLabel)

	if configFilePostfix == "" {
		configFilePostfix = config.FilePostfixDefault
	}

	log.Println("Load ConfigType End")

	log.Println("Load Config Start")
	cfg, err := config.Load(serviceType, configDir, configFilePostfix)

	if err != nil {
		log.Fatal("Load config err: ", err)
		return
	}

	log.Println("Load Config End")

	log.Println("Load Log Start")

	err = logger.Init(serviceType, cfg.Log.Level, cfg.Log.FilePath)

	if err != nil {
		log.Fatal("Init logger err: ", err)
		return
	}

	defer func() {
		err = logger.Sync()

		if err != nil {
			logger.Fatal("Close Logger Fail")
		}
	}()

	logger.Info("Init Log End")

	logger.Info("Init Database Start")

	db, err := database.NewPostgresDB(&database.PostgresConfig{
		Host:         cfg.Database.Host,
		Port:         cfg.Database.Port,
		User:         cfg.Database.User,
		Password:     cfg.Database.Password,
		DBName:       cfg.Database.DBName,
		SSLMode:      cfg.Database.SSLMode,
		MaxIdleConns: cfg.Database.MaxIdleConns,
		MaxOpenConns: cfg.Database.MaxOpenConns,
		MaxLifeTime:  cfg.Database.MaxLifeTime,
		LogLevel:     cfg.Database.LogLevel,
	})

	if err != nil {
		logger.Fatal("Init Database Fail", zap.Error(err))
		return
	}

	defer func() {
		err := database.Close()
		if err != nil {
			logger.Fatal("Close Database Fail", zap.Error(err))
			return
		}
	}()

	if cfg.Database.AutoMigrate {
		logger.Info("AutoMigrate Start")
		err := db.AutoMigrate(&model.User{}, &model.Habit{})

		if err != nil {
			logger.Fatal("AutoMigrate Fail", zap.Error(err))
			return
		}
		logger.Info("AutoMigrate End")
	}
	logger.Info("Init Database End")

	logger.Info("Init Utils Start")
	utils.InitJwt(cfg.Session.Secret)
	logger.Info("Init Utils End")

	logger.Info("Init Redis Start")
	redis, err := database.NewRedis(&database.RedisConfig{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
		MaxRetries:   cfg.Redis.MaxRetries,
		DialTimeout:  time.Duration(cfg.Redis.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.Redis.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Redis.WriteTimeout) * time.Second,
	})

	if err != nil {
		logger.Fatal("Init Redis Fail", zap.Error(err))
		return
	}

	defer func() {
		err := database.CloseRedis()

		if err != nil {
			logger.Fatal("Close Redis Fail", zap.Error(err))
			return
		}
	}()
	logger.Info("Init Redis End")

	logger.Info("Init Repo Start")
	healthRepo := repository.NewHealthRepository()
	userRepo := repository.NewUserRepository(db)
	habitRepo := repository.NewHabitRepository(db)
	logger.Info("Init Repo End")

	logger.Info("Init Service Start")
	healthService := service.NewHealthService(healthRepo)
	userService := service.NewUserService(userRepo, habitRepo)
	habitService := service.NewHabitService(habitRepo, userRepo)
	authService := service.NewAuthService(userRepo, redis)
	logger.Info("Init Service End")

	logger.Info("Init Controller Start")
	healthController := controller.NewHealthController(healthService)
	userController := controller.NewUserController(userService)
	habitController := controller.NewHabitsController(habitService)
	authController := controller.NewAuthController(authService)
	logger.Info("Init Controller End")

	logger.Info("Setup Router Start")
	r := router.SetupRouter(cfg, redis, healthController, userController, habitController, authController)

	if r == nil {
		logger.Fatal("New router err")
		return
	}

	logger.Info("Setup Router End")

	logger.Info("Init Router End")

	logger.Info("Prepare Http Server Start")
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}
	logger.Info("Prepare Http Server End\n\n")

	logger.Info("------- Server Info -------")
	logger.Info(fmt.Sprintf("Local URL: 127.0.0.1:%d", cfg.Server.Port))
	logger.Info(fmt.Sprintf("Test URL Health: 127.0.0.1:%d/health", cfg.Server.Port))
	logger.Info("---------------------------\n\n")

	logger.Info("Start Http Server Start")
	go func() {
		err = server.ListenAndServe()

		if err != nil {
			logger.Fatal("ListenAndServe: ", zap.Error(err))
		}
	}()
	logger.Info("Start Http Server End")

	logger.Info("Listen quit signal ....")
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutdown Server ...")

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(cfg.Server.CloseTimeout)*time.Second)
	defer cancelFunc()

	err = server.Shutdown(ctx)

	if err != nil {
		logger.Fatal("Server Shutdown Error: ", zap.Error(err))
		return
	}

	logger.Info("Server exiting")
}
