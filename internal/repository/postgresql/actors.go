package postgresql

import (
	"context"
	"github.com/google/uuid"
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

func (r *ActorsRepo) Create(ctx context.Context, actor models.Actor) (uuid.UUID, error) {
	var id uuid.UUID

	query := `INSERT INTO actors (fio, birthday, sex) VALUES (row(@name, @secondName, @patronymic), @bd, @s) RETURNING id`
	args := pgx.NamedArgs{
		"name":       actor.Name,
		"secondName": actor.SecondName,
		"patronymic": actor.Patronymic,
		"bd":         actor.DateOfBirth,
		"s":          actor.Sex,
	}

	err := r.db.QueryRow(ctx, query, args).Scan(&id)

	return id, err
}

func (r *ActorsRepo) InsertIntoActorFilm(ctx context.Context, actorId uuid.UUID, filmId uuid.UUID) error {
	query := `INSERT INTO actors_films (fk_actor_id, fk_film_id) VALUES (@actor, @film)`
	args := pgx.NamedArgs{
		"actor": actorId,
		"film":  filmId,
	}

	_, err := r.db.Exec(ctx, query, args)

	return err
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
