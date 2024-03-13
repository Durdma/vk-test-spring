package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
	"vk-test-spring/internal/config"
	"vk-test-spring/pkg/logger"
)

const timeout = 10 * time.Second

func NewConnection(cfg config.PostgreSQLConfig) *pgx.Conn {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := pgx.Connect(ctx, getConnectionString(cfg))
	if err != nil {
		logger.Error(err.Error())
	}

	err = conn.Ping(context.Background())
	if err != nil {
		logger.Error(err.Error())
	}

	return conn
}

func getConnectionString(cfg config.PostgreSQLConfig) string {
	return fmt.Sprintf("postgresql://%v:%v@%v:%v/%v", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
}
