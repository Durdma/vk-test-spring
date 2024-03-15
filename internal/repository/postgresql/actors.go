package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"vk-test-spring/internal/models"
	"vk-test-spring/pkg/logger"
)

type ActorsRepo struct {
	db *pgxpool.Pool
}

func MewActorsRepo(db *pgxpool.Pool) *ActorsRepo {
	return &ActorsRepo{
		db: db,
	}
}

func (r *ActorsRepo) Create(ctx context.Context, actor models.Actor) (uuid.UUID, error) {
	var id uuid.UUID

	query := `INSERT INTO actors (f_name, s_name, patronymic, birthday, sex) VALUES (@name, @secondName, @patron, @bd, @s) RETURNING id`
	args := pgx.NamedArgs{
		"name":       actor.Name,
		"secondName": actor.SecondName,
		"patron":     actor.Patronymic,
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

func (r *ActorsRepo) Delete(ctx context.Context, actorId uuid.UUID) error {
	query := `DELETE FROM actors WHERE id=@actorId`
	args := pgx.NamedArgs{
		"actorId": actorId,
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
	return nil
}

func (r *ActorsRepo) GetAllActors(ctx context.Context) ([]models.Actor, error) {
	query := `SELECT id, f_name, s_name, patronymic, birthday, sex FROM actors`

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		tx.Rollback(ctx)
		return nil, err
	}
	defer rows.Close()

	actors := make([]models.Actor, 0)

	for rows.Next() {
		actor := models.Actor{}
		var t time.Time

		err := rows.Scan(&actor.ID, &actor.Name, &actor.SecondName, &actor.Patronymic, &t, &actor.Sex)
		if err != nil {
			tx.Rollback(ctx)
			return nil, err
		}

		actor.DateOfBirth = r.dateTypeToString(t)

		films, err := r.getActorFilms(ctx, actor.ID)
		if err != nil {

			tx.Rollback(ctx)
			return nil, err
		}

		actor.Films = films

		fmt.Println(actor)

		actors = append(actors, actor)
	}

	tx.Commit(ctx)
	return actors, nil
}

func (r *ActorsRepo) GetActorsByName(ctx context.Context, name string) ([]models.Actor, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, `SELECT actors.id, actors.f_name, actors.s_name, actors.patronymic, actors.birthday, actors.sex
	FROM actors WHERE concat(actors.f_name, ' ', actors.s_name, ' ', actors.patronymic) LIKE '%$1%'`, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		tx.Rollback(ctx)
		return nil, err
	}
	defer rows.Close()

	actors := make([]models.Actor, 0)

	for rows.Next() {
		actor := models.Actor{}
		var t time.Time

		err := rows.Scan(&actor.ID, &actor.Name, &actor.SecondName, &actor.Patronymic, &t, &actor.Sex)
		if err != nil {
			tx.Rollback(ctx)
			return nil, err
		}

		actor.DateOfBirth = r.dateTypeToString(t)

		films, err := r.getActorFilms(ctx, actor.ID)
		if err != nil {

			tx.Rollback(ctx)
			return nil, err
		}

		actor.Films = films

		fmt.Println(actor)

		actors = append(actors, actor)
	}

	tx.Commit(ctx)
	return actors, nil
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

func (r *ActorsRepo) getActorFilms(ctx context.Context, actorId uuid.UUID) ([]models.ActorFilm, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, `SELECT films.id, films.name
	FROM films
	JOIN actors_films ON films.id = actors_films.fk_film_id
	JOIN actors ON actors.id = actors_films.fk_actor_id
	WHERE actors.id = $1`, actorId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		logger.Errorf("error 1 in getActorFilms: %v", err)

		tx.Rollback(ctx)
		return nil, err
	}
	defer rows.Close()

	var films []models.ActorFilm
	for rows.Next() {
		film := models.ActorFilm{}

		err := rows.Scan(&film.ID, &film.Name)
		if err != nil {
			logger.Errorf("error 2 in getActorFilms: %v", err)

			tx.Rollback(ctx)
			return nil, err
		}

		films = append(films, film)
	}

	tx.Commit(ctx)
	return films, nil
}
