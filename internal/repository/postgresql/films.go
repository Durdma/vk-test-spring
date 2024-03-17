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

type FilmsRepo struct {
	db *pgxpool.Pool
}

func NewFilmsRepo(db *pgxpool.Pool) *FilmsRepo {
	return &FilmsRepo{
		db: db,
	}
}

func (r *FilmsRepo) Create(ctx context.Context, film models.Film, actors []uuid.UUID) error {
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
		return err
	}

	err = tx.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	err = r.insertIntoActorFilm(ctx, tx, actors, id)
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

	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	err = r.insertIntoActorFilm(ctx, tx, actorsToAdd, film.ID)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	err = r.deleteFromActorFilm(ctx, tx, actorsToDel, film.ID)
	if err != nil {
		tx.Rollback(ctx)
		return err
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

	res, err := tx.Exec(ctx, query, args)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	if res.RowsAffected() == 0 {
		tx.Rollback(ctx)
		return models.CustomError{Code: http.StatusNotFound, Message: fmt.Sprintf("not found film with this id: %v", filmId)}
	}

	tx.Commit(ctx)
	return err
}

func (r *FilmsRepo) GetAllFilms(ctx context.Context) ([]models.Film, error) {
	query := fmt.Sprintf("SELECT id, name, description, date, rating FROM films ORDER BY %s %s", ctx.Value("sort"), ctx.Value("order"))

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			tx.Commit(ctx)
			return nil, models.CustomError{Code: http.StatusNotFound, Message: fmt.Sprintln("films not found")}
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

		actors, err := r.getFilmActors(ctx, tx, film.ID)
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

func (r *FilmsRepo) GetFilmByName(ctx context.Context, name string) ([]models.Film, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, fmt.Sprintf(`SELECT id, name, description, date, rating FROM films WHERE films.name LIKE '%%%v%%'`, name))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			tx.Commit(ctx)
			return nil, models.CustomError{Code: http.StatusNotFound, Message: fmt.Sprintf("not found film with name like this: %v", name)}
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

		actors, err := r.getFilmActors(ctx, tx, film.ID)
		if err != nil {
			tx.Rollback(ctx)
			return nil, err
		}

		film.Actors = actors

		films = append(films, film)
	}

	tx.Commit(ctx)
	return films, nil
}

func (r *FilmsRepo) GetFilmByActor(ctx context.Context, actorName string) ([]models.Film, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, fmt.Sprintf("SELECT films.id, films.description, films.date, films.rating"+
		" FROM films JOIN actors_films as af ON films.id = af.fk_film_id JOIN actors as a ON af.fk_actor_id = a.id"+
		" WHERE CONCAT(a.f_name, ' ', a.s_name, ' ', a.patronymic) LIKE '%%%v%%'", actorName))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			tx.Commit(ctx)
			return nil, models.CustomError{Code: http.StatusNotFound, Message: fmt.Sprintf("not found film with actor name like this: %v", actorName)}
		}

		tx.Rollback(ctx)
		return nil, err
	}
	defer rows.Close()

	films := make([]models.Film, 0)
	for rows.Next() {
		film := models.Film{}
		var t time.Time

		err := rows.Scan(&film.ID, &film.Description, &t, &film.Rating)
		if err != nil {
			tx.Rollback(ctx)
			return nil, err
		}

		film.Date = r.dateTypeToString(t)

		actors, err := r.getFilmActors(ctx, tx, film.ID)
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

	err = tx.QueryRow(ctx, `SELECT id, name, description, date, rating
	FROM films WHERE id=$1`, filmId).Scan(&film.ID, &film.Name, &film.Description, &t, &film.Rating)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			tx.Rollback(ctx)
			return models.Film{}, models.CustomError{Code: http.StatusNotFound, Message: fmt.Sprintf("not found film with this id: %v", filmId)}
		}
		tx.Rollback(ctx)
		return models.Film{}, err
	}

	film.Date = r.dateTypeToString(t)

	actors, err := r.getFilmActors(ctx, tx, filmId)
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

func (r *FilmsRepo) insertIntoActorFilm(ctx context.Context, tx pgx.Tx, actorsId []uuid.UUID, filmId uuid.UUID) error {
	query := `INSERT INTO actors_films (fk_actor_id, fk_film_id) VALUES (@actor, @film)`
	if len(actorsId) > 0 {
		for _, a := range actorsId {
			args := pgx.NamedArgs{
				"actor": a,
				"film":  filmId,
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

func (r *FilmsRepo) deleteFromActorFilm(ctx context.Context, tx pgx.Tx, actorsId []uuid.UUID, filmId uuid.UUID) error {
	query := `DELETE FROM actors_films WHERE fk_actor_id = $1 AND fk_film_id = $2`
	if len(actorsId) > 0 {
		for _, a := range actorsId {
			_, err := tx.Exec(ctx, query, a, filmId)
			if err != nil {
				return err
			}
		}

		return nil
	}

	return nil
}

func (r *FilmsRepo) getFilmActors(ctx context.Context, tx pgx.Tx, filmId uuid.UUID) ([]models.FilmActors, error) {
	rows, err := r.db.Query(ctx, `SELECT actors.id, actors.f_name, actors.s_name, actors.patronymic
	FROM actors
	JOIN actors_films ON actors.id = actors_films.fk_actor_id
	JOIN films ON films.id = actors_films.fk_film_id
	WHERE films.id = $1`, filmId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}

		return nil, err
	}
	defer rows.Close()

	var actors []models.FilmActors
	for rows.Next() {
		actor := models.FilmActors{}

		err := rows.Scan(&actor.ID, &actor.Name, &actor.SecondName, &actor.Patronymic)
		if err != nil {
			return nil, err
		}

		actors = append(actors, actor)
	}

	return actors, err
}

func (r *FilmsRepo) dateTypeToString(t time.Time) string {
	return t.Format(time.DateOnly)
}
