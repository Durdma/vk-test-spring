package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"vk-test-spring/internal/models"
	"vk-test-spring/internal/repository/postgresql"
)

type Films interface {
	Create(ctx context.Context, film models.Film) (uuid.UUID, error)
	Update(ctx context.Context, film models.Film, actorsToAdd []uuid.UUID, actorsToDel []uuid.UUID) error
	Delete(ctx context.Context, filmId uuid.UUID) error
	GetFilmByName(ctx context.Context, name string) ([]models.Film, error)
	GetFilmByActor(ctx context.Context, actorName string) ([]models.Film, error)
	GetAllFilms(ctx context.Context) ([]models.Film, error)
	InsertIntoActorFilm(ctx context.Context, actorId uuid.UUID, filmId uuid.UUID) error
	GetFilmById(ctx context.Context, filmId uuid.UUID) (models.Film, error)
	DeleteFromActorFilm(ctx context.Context, filmId uuid.UUID, actorId uuid.UUID) error
}

type Actors interface {
	Create(ctx context.Context, actor models.Actor, actorFilms []uuid.UUID) error
	Edit(ctx context.Context, actor models.Actor, filmsToAdd []uuid.UUID, filmsToDel []uuid.UUID) error
	Delete(ctx context.Context, actorId uuid.UUID) error
	GetAllActors(ctx context.Context) ([]models.Actor, error)
	GetActorsByName(ctx context.Context, name string) ([]models.Actor, error)
	GetActorById(ctx context.Context, actorId uuid.UUID) (models.Actor, error)
}

type Users interface {
	Create(ctx context.Context, user models.User) error
	Delete(ctx context.Context, userId string) error
	Edit(ctx context.Context, user models.User) error
	GetUserIdRole(username string, password string) (string, string, error)
}

type Repositories struct {
	Films  Films
	Actors Actors
	Users  Users
}

func NewRepositories(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		Films:  postgresql.NewFilmsRepo(db),
		Actors: postgresql.MewActorsRepo(db),
		Users:  postgresql.NewUsersRepo(db),
	}
}

//func NewRepositories(db *pgx.Conn) *Repositories {
//	return &Repositories{
//		Films:  postgresql.NewFilmsRepo(db),
//		Actors: postgresql.MewActorsRepo(db),
//		Users:  postgresql.NewUsersRepo(db),
//	}
//}
