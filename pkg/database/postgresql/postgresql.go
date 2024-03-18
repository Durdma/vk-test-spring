package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"vk-test-spring/internal/config"
)

const timeout = 10 * time.Second

func NewConnectionPool(cfg config.PostgreSQLConfig) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	pool, err := pgxpool.New(ctx, getConnectionString(cfg))
	if err != nil {
	}

	err = pool.Ping(context.Background())
	if err != nil {
	}

	return pool
}

func NewConnection(cfg config.PostgreSQLConfig) *pgx.Conn {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := pgx.Connect(ctx, getConnectionString(cfg))
	if err != nil {
	}

	err = conn.Ping(context.Background())
	if err != nil {
	}

	return conn
}

func getConnectionString(cfg config.PostgreSQLConfig) string {
	return fmt.Sprintf("postgresql://%v:%v@%v:%v/%v", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
}
