package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"time"
	"vk-test-spring/internal/models"
)

type ActorsRepo struct {
	db *pgxpool.Pool
}

func MewActorsRepo(db *pgxpool.Pool) *ActorsRepo {
	return &ActorsRepo{
		db: db,
	}
}

func (r *ActorsRepo) Create(ctx context.Context, actor models.Actor, actorFilms []uuid.UUID) error {
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
		return err
	}

	err = tx.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	err = r.insertIntoActorFilms(ctx, tx, id, actorFilms)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	tx.Commit(ctx)
	return err
}

func (r *ActorsRepo) Edit(ctx context.Context, actor models.Actor, filmsToAdd []uuid.UUID, filmsToDel []uuid.UUID) error {
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

	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	err = r.insertIntoActorFilms(ctx, tx, actor.ID, filmsToAdd)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	err = r.deleteFromActorFilms(ctx, tx, actor.ID, filmsToDel)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	tx.Commit(ctx)
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

	res, err := tx.Exec(ctx, query, args)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	if res.RowsAffected() == 0 {
		tx.Rollback(ctx)
		return models.CustomError{Code: http.StatusNotFound, Message: fmt.Sprintf("not found actor with this id: %v", actorId)}
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

	rows, err := tx.Query(ctx, query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			tx.Commit(ctx)
			return nil, models.CustomError{Code: http.StatusNotFound, Message: fmt.Sprintln("actors not found")}
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

		films, err := r.getActorFilms(ctx, tx, actor.ID)
		if err != nil {
			tx.Rollback(ctx)
			return nil, err
		}

		actor.Films = films

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

	rows, err := tx.Query(ctx, `SELECT actors.id, actors.f_name, actors.s_name, actors.patronymic, actors.birthday, actors.sex
	FROM actors WHERE concat(actors.f_name, ' ', actors.s_name, ' ', actors.patronymic) LIKE '%$1%'`, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			tx.Commit(ctx)
			return nil, models.CustomError{Code: http.StatusNotFound, Message: fmt.Sprintf("not found. actors with name like this: %v", name)}
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

		films, err := r.getActorFilms(ctx, tx, actor.ID)
		if err != nil {
			tx.Rollback(ctx)
			return nil, err
		}

		actor.Films = films

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

	err = tx.QueryRow(ctx, `SELECT id, f_name, s_name, patronymic, birthday, sex 
	FROM actors WHERE id=$1`, actorId).Scan(
		&actor.ID, &actor.Name, &actor.SecondName, &actor.Patronymic, &t, &actor.Sex)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Actor{}, models.CustomError{Code: http.StatusNotFound, Message: fmt.Sprintf("not found actor with this id: %v", actorId)}
		}
		tx.Rollback(ctx)
		return models.Actor{}, err
	}

	actor.DateOfBirth = r.dateTypeToString(t)

	films, err := r.getActorFilms(ctx, tx, actorId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			tx.Commit(ctx)
			return actor, nil
		}

		tx.Rollback(ctx)
		return models.Actor{}, err
	}

	actor.Films = films

	tx.Commit(ctx)
	return actor, err
}

func (r *ActorsRepo) insertIntoActorFilms(ctx context.Context, tx pgx.Tx, actorId uuid.UUID, filmsId []uuid.UUID) error {
	query := `INSERT INTO actors_films (fk_actor_id, fk_film_id) VALUES (@actor, @film)`
	if len(filmsId) > 0 {
		for _, f := range filmsId {
			args := pgx.NamedArgs{
				"actor": actorId,
				"film":  f,
			}

			_, err := tx.Exec(ctx, query, args)
			if err != nil {
				return err
			}
		}

		return nil
	}

	return nil
}

func (r *ActorsRepo) deleteFromActorFilms(ctx context.Context, tx pgx.Tx, actorId uuid.UUID, filmsId []uuid.UUID) error {
	query := `DELETE FROM actors_films WHERE fk_actor_id = $1 AND fk_film_id = $2`
	if len(filmsId) > 0 {
		for _, f := range filmsId {
			_, err := tx.Exec(ctx, query, actorId, f)
			if err != nil {
				return err
			}
		}

		return nil
	}

	return nil
}

func (r *ActorsRepo) getActorFilms(ctx context.Context, tx pgx.Tx, actorId uuid.UUID) ([]models.ActorFilm, error) {
	rows, err := r.db.Query(ctx, `SELECT films.id, films.name
	FROM films
	JOIN actors_films ON films.id = actors_films.fk_film_id
	JOIN actors ON actors.id = actors_films.fk_actor_id
	WHERE actors.id = $1`, actorId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}
	defer rows.Close()

	var films []models.ActorFilm
	for rows.Next() {
		film := models.ActorFilm{}

		err := rows.Scan(&film.ID, &film.Name)
		if err != nil {
			return nil, err
		}

		films = append(films, film)
	}

	return films, err
}

func (r *ActorsRepo) dateTypeToString(t time.Time) string {
	return t.Format(time.DateOnly)
}
