package postgresql

import (
	"context"
	"github.com/jackc/pgx/v5"
	"vk-test-spring/internal/models"
)

type ActorsRepo struct {
	db *pgx.Conn
}

func MewActorsRepo(db *pgx.Conn) *ActorsRepo {
	return &ActorsRepo{
		db: db,
	}
}

func (r *ActorsRepo) Create(ctx context.Context, actor models.Actor) error {
	return nil
}

func (r *ActorsRepo) Edit(ctx context.Context, actor models.Actor) error {
	return nil
}

func (r *ActorsRepo) Delete(ctx context.Context, actorId string) error {
	return nil
}

func (r *ActorsRepo) GetAllActors(ctx context.Context) ([]models.Actor, error) {
	return nil, nil
}

func (r *ActorsRepo) GetActorsByName(ctx context.Context, name string) ([]models.Actor, error) {
	return nil, nil
}
