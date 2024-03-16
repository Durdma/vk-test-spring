package postgresql

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
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

func (r *FilmsRepo) InsertIntoActorFilm(ctx context.Context, actorId uuid.UUID, filmId uuid.UUID) error {
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

func (r *FilmsRepo) DeleteFromActorFilm(ctx context.Context, filmId uuid.UUID, actorId uuid.UUID) error {
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

func (r *FilmsRepo) Update(ctx context.Context, film models.Film, actorsToAdd []uuid.UUID, actorsToDel []uuid.UUID) error {
	query := `UPDATE films SET name = @n, description = @d, date = @dd, rating = @r WHERE id=@film_id`
	args := pgx.NamedArgs{
		"n":       film.Name,
		"d":       film.Description,
		"dd":      film.Date,
		"r":       film.Rating,
		"film_id": film.ID,
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

	if len(actorsToAdd) > 0 {
		for _, a := range actorsToAdd {
			err := r.InsertIntoActorFilm(ctx, a, film.ID)
			if err != nil {
				tx.Rollback(ctx)
				return err
			}
		}
	}

	if len(actorsToDel) > 0 {
		for _, a := range actorsToDel {
			err := r.DeleteFromActorFilm(ctx, film.ID, a)
			if err != nil {
				tx.Rollback(ctx)
				return err
			}
		}
	}

	tx.Commit(ctx)
	return err
}

func (r *FilmsRepo) Delete(ctx context.Context, filmId uuid.UUID) error {
	query := `DELETE FROM films WHERE id=@filmId`
	args := pgx.NamedArgs{
		"filmId": filmId,
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

func (r *FilmsRepo) GetFilmByName(ctx context.Context, name string) ([]models.Film, error) {
	return nil, nil
}

func (r *FilmsRepo) GetFilmByActor(ctx context.Context, actorName string) ([]models.Film, error) {
	return nil, nil
}

func (r *FilmsRepo) GetAllFilms(ctx context.Context) ([]models.Film, error) {
	query := `SELECT id, name, description, date, rating FROM films`

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

	films := make([]models.Film, 0)
	for rows.Next() {
		film := models.Film{}
		var t time.Time

		err := rows.Scan(&film.ID, &film.Name, &film.Description, &t, &film.Rating)
		if err != nil {
			tx.Rollback(ctx)
			return nil, err
		}

		film.Date = r.dateTypeToString(t)

		actors, err := r.getFilmActors(ctx, film.ID)
		if err != nil {
			tx.Rollback(ctx)
			return nil, err
		}

		film.Actors = actors

		films = append(films, film)
	}

	tx.Commit(ctx)
	return films, err
}

func (r *FilmsRepo) GetFilmById(ctx context.Context, filmId uuid.UUID) (models.Film, error) {
	var film models.Film
	var t time.Time

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return models.Film{}, err
	}

	err = r.db.QueryRow(ctx, `SELECT id, name, description, date, rating
	FROM films WHERE id=$1`, filmId).Scan(&film.ID, &film.Name, &film.Description, &t, &film.Rating)
	if err != nil {
		tx.Rollback(ctx)
		return models.Film{}, err
	}

	film.Date = r.dateTypeToString(t)

	actors, err := r.getFilmActors(ctx, filmId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			tx.Commit(ctx)
			return film, nil
		}

		tx.Rollback(ctx)
		return models.Film{}, err
	}

	film.Actors = actors

	tx.Commit(ctx)
	return film, err
}

func (r *FilmsRepo) getFilmActors(ctx context.Context, filmId uuid.UUID) ([]models.FilmActors, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, `SELECT actors.id, actors.f_name, actors.s_name, actors.patronymic
	FROM actors
	JOIN actors_films ON actors.id = actors_films.fk_actor_id
	JOIN films ON films.id = actors_films.fk_film_id
	WHERE films.id = $1`, filmId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {

			tx.Commit(ctx)
			return nil, nil
		}

		tx.Rollback(ctx)
		return nil, err
	}
	defer rows.Close()

	var actors []models.FilmActors
	for rows.Next() {
		actor := models.FilmActors{}

		err := rows.Scan(&actor.ID, &actor.Name, &actor.SecondName, &actor.Patronymic)
		if err != nil {

			tx.Rollback(ctx)
			return nil, err
		}

		actors = append(actors, actor)
	}

	tx.Commit(ctx)
	return actors, err
}

func (r *FilmsRepo) dateTypeToString(t time.Time) string {
	return t.Format(time.DateOnly)
}
