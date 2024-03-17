package service

import (
	"context"
	"github.com/google/uuid"
	"vk-test-spring/internal/models"
	"vk-test-spring/internal/repository"
)

type Films interface {
	AddNewFilm(ctx context.Context, input FilmCreateInput) error
	EditFilm(ctx context.Context, input FilmUpdateInput) error
	DeleteFilm(ctx context.Context, filmId uuid.UUID) error
	GetAllFilms(ctx context.Context) ([]models.Film, error)
	GetAllFilmsByName(ctx context.Context, name string) ([]models.Film, error)
	GetAllFilmsByActor(ctx context.Context, actorsName string) ([]models.Film, error)
}

type Actors interface {
	AddActor(ctx context.Context, input ActorCreateInput) error
	UpdateActor(ctx context.Context, input ActorUpdateInput) error
	DeleteActor(ctx context.Context, actorId uuid.UUID) error
	GetAllActors(ctx context.Context) ([]models.Actor, error)
	GetActorById(ctx context.Context, actorId uuid.UUID) (models.Actor, error)
	GetActorByName(ctx context.Context, name string) ([]models.Actor, error)
}

type Users interface {
	CreateUser(ctx context.Context, input UserInput) error
	DeleteUser(ctx context.Context, userId string) error
	ChangeRole(ctx context.Context, userId string, role string) error
	GetUserIdRole(username string, password string) (string, string, error)
}

type Services struct {
	Films  Films
	Actors Actors
	Users  Users
}

func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		Films:  NewFilmsService(repos.Films),
		Actors: NewActorsService(repos.Actors),
		Users:  NewUsersService(repos.Users),
	}
}
