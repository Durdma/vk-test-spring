package postgresql

import (
	"context"
	"github.com/jackc/pgx/v5"
	"vk-test-spring/internal/models"
)

type FilmsRepo struct {
	db *pgx.Conn
}

func NewFilmsRepo(db *pgx.Conn) *FilmsRepo {
	return &FilmsRepo{
		db: db,
	}
}

func (r *FilmsRepo) Create(ctx context.Context, film models.Film) error {
	return nil
}

func (r *FilmsRepo) Update(ctx context.Context, film models.Film) error {
	return nil
}

func (r *FilmsRepo) Delete(ctx context.Context, filmId string) error {
	return nil
}

func (r *FilmsRepo) GetFilmByName(ctx context.Context, name string) ([]models.Film, error) {
	return nil, nil
}

func (r *FilmsRepo) GetFilmByActor(ctx context.Context, actorName string) ([]models.Film, error) {
	return nil, nil
}

func (r *FilmsRepo) GetAllFilms(ctx context.Context) ([]models.Film, error) {
	return nil, nil
}
