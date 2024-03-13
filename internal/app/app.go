package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"vk-test-spring/internal/config"
	"vk-test-spring/internal/server"
	"vk-test-spring/pkg/logger"
)

func Run(configPath string) {
	cfg, err := config.Init(configPath)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Infof("%+v\n", *cfg)

	srv := server.NewServer(cfg)
	go func() {
		if err := srv.Run(); err != nil {
			fmt.Printf("error while running http server: %s\n", err.Error())
		}
	}()

	logger.Info("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
}
