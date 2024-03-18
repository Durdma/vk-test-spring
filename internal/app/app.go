package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"vk-test-spring/internal/config"
	"vk-test-spring/internal/controller"
	"vk-test-spring/internal/repository"
	"vk-test-spring/internal/server"
	"vk-test-spring/internal/service"
	"vk-test-spring/pkg/database/postgresql"
	"vk-test-spring/pkg/logger"
)

// @title Films library API
// @version 1.0
// @description API Server for Films library

// @host localhost:8080
// @BasePath /films

// @SecurityDefinitions basicAuth
// @SecurityScheme basic
// @Security BasicAuth

func Run(configPath string) {
	// TODO add config for logger
	logs := logger.InitLogs("../../pkg/logger/logger.json", 5, 3, 30)

	logs.Info().Msg("Starting app")

	cfg, err := config.Init(configPath)
	if err != nil {
		return
	}

	dbHandler := postgresql.NewConnectionPool(cfg.PostgreSQL)
	logs.Info().Msg("Initialized connection pool DB")

	repos := repository.NewRepositories(dbHandler)
	logs.Info().Msg("Initialized repos")

	services := service.NewServices(repos)
	logs.Info().Msg("Initialized services")

	handlers := controller.NewHandler()
	mux := handlers.Init(services, logs)
	logs.Info().Msg("Initialized handlers")

	srv := server.NewServer(cfg, mux)
	go func() {
		if err := srv.Run(); err != nil {
			logs.Error().Msg(fmt.Sprintf("error while starting server: %v", err.Error()))
		}
	}()

	logs.Info().Msg("server started")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	dbHandler.Close()

	//if err := dbHandler.Close(context.Background()); err != nil {
	//	logger.Error(err.Error())
	//}
}
