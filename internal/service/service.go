package service

import (
	"context"
	"vk-test-spring/internal/models"
	"vk-test-spring/internal/repository"
)

type FilmInput struct {
}

type Films interface {
	AddNewFilm(ctx context.Context, input FilmInput) error
	EditFilm(ctx context.Context, input FilmInput) error
	DeleteFilm(ctx context.Context, name string) error
	GetAllFilms(ctx context.Context) ([]models.Film, error)
	GetAllFilmsByName(ctx context.Context, name string) ([]models.Film, error)
	GetAllFilmsByActor(ctx context.Context, actorsName string) ([]models.Film, error)
}

type ActorInput struct {
}

type Actors interface {
	AddActor(ctx context.Context, input ActorInput) error
	UpdateActor(ctx context.Context, actor models.Actor) error
	DeleteActor(ctx context.Context, actorId string) error
	GetAllActors(ctx context.Context) ([]models.Actor, error)
}

type UserInput struct {
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
