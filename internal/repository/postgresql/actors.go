package postgresql

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"time"
	"vk-test-spring/internal/models"
	"vk-test-spring/pkg/logger"
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

	query := `INSERT INTO actors (f_name, s_name, patronymic, birthday, sex) VALUES (@name, @secondName, @patronymic, @bd, @s) RETURNING id`
	args := pgx.NamedArgs{
		"name":       actor.Name,
		"secondName": actor.SecondName,
		"patronymic": actor.Patronymic,
		"bd":         actor.DateOfBirth,
		"s":          actor.Sex,
	}

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return uuid.Nil, err
	}

	err = r.db.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		tx.Rollback(ctx)
		return uuid.Nil, err
	}

	tx.Commit(ctx)
	return id, err
}

func (r *ActorsRepo) InsertIntoActorFilm(ctx context.Context, actorId uuid.UUID, filmId uuid.UUID) error {
	query := `INSERT INTO actors_films (fk_actor_id, fk_film_id) VALUES (@actor, @film)`
	args := pgx.NamedArgs{
		"actor": actorId,
		"film":  filmId,
	}

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	tx.Commit(ctx)
	return err
}

func (r *ActorsRepo) DeleteFromActorFilm(ctx context.Context, actorId uuid.UUID, filmId uuid.UUID) error {
	query := `DELETE FROM actors_films WHERE fk_actor_id = @actorId AND fk_film_id = @filmId`
	args := pgx.NamedArgs{
		"actorId": actorId,
		"filmId":  filmId,
	}

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	tx.Commit(ctx)
	return err
}

func (r *ActorsRepo) Edit(ctx context.Context, actor models.Actor) error {
	query := `UPDATE actors SET f_name = @name, s_name = @secondName, patronymic = @patronymic, birthday = @bd, sex = @s WHERE id = @actor_id`
	args := pgx.NamedArgs{
		"name":       actor.Name,
		"secondName": actor.SecondName,
		"patronymic": actor.Patronymic,
		"bd":         actor.DateOfBirth,
		"s":          actor.Sex,
		"actor_id":   actor.ID,
	}

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	logger.Errorf("error in Edit: %v", err)

	return err
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

func (r *ActorsRepo) GetActorById(ctx context.Context, actorId uuid.UUID) (models.Actor, error) {
	var actor models.Actor
	var t time.Time

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return models.Actor{}, err
	}

	err = r.db.QueryRow(ctx, "SELECT id, f_name, s_name, patronymic, birthday, sex "+
		"FROM actors WHERE id=$1", actorId).Scan(
		&actor.ID, &actor.Name, &actor.SecondName, &actor.Patronymic, &t, &actor.Sex)
	if err != nil {
		logger.Errorf("error 1 in GetActorById: %v", err)

		tx.Rollback(ctx)
		return models.Actor{}, err
	}

	actor.DateOfBirth = r.dateTypeToString(t)

	films, err := r.getActorFilms(ctx, actorId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return actor, nil
		}

		logger.Errorf("error 2 in GetActorById: %v", err)
		tx.Rollback(ctx)
		return models.Actor{}, err
	}

	actor.Films = films

	tx.Commit(ctx)
	return actor, err
}

func (r *ActorsRepo) dateTypeToString(t time.Time) string {
	return t.Format(time.DateOnly)
}

func (r *ActorsRepo) getActorFilms(ctx context.Context, actorId uuid.UUID) ([]models.Film, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, `SELECT films.id, films.name, films.description, films.date, films.rating
	FROM films
	JOIN actors_films ON films.id = actors_films.fk_film_id
	JOIN actors ON actors.id = actors_films.fk_actor_id
	WHERE actors.id = $1`, actorId)
	if err != nil {
		logger.Errorf("error 1 in getActorFilms: %v", err)

		tx.Rollback(ctx)
		return nil, err
	}
	defer rows.Close()

	var films []models.Film
	for rows.Next() {
		film := models.Film{}
		var t time.Time

		err := rows.Scan(&film.ID, &film.Name, &film.Description, &t, &film.Rating)
		if err != nil {
			logger.Errorf("error 2 in getActorFilms: %v", err)

			tx.Rollback(ctx)
			return nil, err
		}

		film.Date = r.dateTypeToString(t)

		films = append(films, film)
	}

	tx.Commit(ctx)
	return films, nil
}
