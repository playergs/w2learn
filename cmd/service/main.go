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
	"w2learn/internal/respository"
	"w2learn/internal/router"
	"w2learn/internal/service"
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

	logger.Info("Init Router Start")
	logger.Info("Init Health Router Start")

	healthRepo := respository.NewHealthRepository()
	healthService := service.NewHealthService(healthRepo)
	healthController := controller.NewHealthController(healthService)

	logger.Info("Init Health Router End")

	logger.Info("Setup Router Start")
	r := router.SetupRouter(cfg, healthController)

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
