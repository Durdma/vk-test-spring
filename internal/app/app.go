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

func Run(configPath string) {
	cfg, err := config.Init(configPath)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Infof("%+v\n", *cfg)

	//dbHandler := postgresql.NewConnection(cfg.PostgreSQL)

	dbHandler := postgresql.NewConnectionPool(cfg.PostgreSQL)

	repos := repository.NewRepositories(dbHandler)

	services := service.NewServices(repos)

	handlers := controller.NewHandler()
	mux := handlers.Init(services)

	srv := server.NewServer(cfg, mux)
	go func() {
		if err := srv.Run(); err != nil {
			fmt.Printf("error while running http server: %s\n", err.Error())
		}
	}()

	logger.Info("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	dbHandler.Close()

	//if err := dbHandler.Close(context.Background()); err != nil {
	//	logger.Error(err.Error())
	//}
}
