package postgresql

import (
	"context"
	"github.com/jackc/pgx/v5"
	"vk-test-spring/internal/models"
)

type UsersRepo struct {
	db *pgx.Conn
}

func NewUsersRepo(db *pgx.Conn) *UsersRepo {
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
