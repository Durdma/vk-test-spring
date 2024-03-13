package app

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
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

	connString := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v", cfg.PostgreSQL.User, cfg.PostgreSQL.Password, cfg.PostgreSQL.Host, cfg.PostgreSQL.Port, cfg.PostgreSQL.DBName)

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		logger.Error(err.Error())
	}
	defer conn.Close(context.Background())

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
