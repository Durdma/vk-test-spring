package postgresql

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"vk-test-spring/internal/models"
	"vk-test-spring/pkg/logger"
)

type UsersRepo struct {
	db *pgxpool.Pool
}

func NewUsersRepo(db *pgxpool.Pool) *UsersRepo {
	return &UsersRepo{
		db: db,
	}
}

func (r *UsersRepo) Create(ctx context.Context, user models.User) error {
	return nil
}

func (r *UsersRepo) Delete(ctx context.Context, userId string) error {
	return nil
}

func (r *UsersRepo) Edit(ctx context.Context, user models.User) error {
	return nil
}

func (r *UsersRepo) GetUserIdRole(username string, password string) (string, string, error) {
	userId := ""
	role := ""

	row := r.db.QueryRow(context.Background(), "SELECT id, role FROM users WHERE name=$1 AND password=$2", username, password).Scan(&userId, &role)
	if row != nil {
		logger.Error(row.Error())
	}

	return userId, role, row
}
