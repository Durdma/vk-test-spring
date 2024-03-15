package postgresql

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"vk-test-spring/internal/models"
)

type FilmsRepo struct {
	db *pgxpool.Pool
}

func NewFilmsRepo(db *pgxpool.Pool) *FilmsRepo {
	return &FilmsRepo{
		db: db,
	}
}

func (r *FilmsRepo) Create(ctx context.Context, film models.Film) (uuid.UUID, error) {
	var id uuid.UUID

	query := `INSERT INTO films (name, description, date, rating) VALUES (@n, @des, @d, @rate) RETURNING id`
	args := pgx.NamedArgs{
		"n":    film.Name,
		"des":  film.Description,
		"d":    film.Date,
		"rate": film.Rating,
	}

	err := r.db.QueryRow(ctx, query, args).Scan(&id)

	return id, err
}

func (r *FilmsRepo) InsertIntoActorFilm(ctx context.Context, actorId uuid.UUID, filmId uuid.UUID) error {
	query := `INSERT INTO actors_films (fk_actor_id, fk_film_id) VALUES (@actor, @film)`
	args := pgx.NamedArgs{
		"actor": actorId,
		"film":  filmId,
	}

	_, err := r.db.Exec(ctx, query, args)

	return err
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
