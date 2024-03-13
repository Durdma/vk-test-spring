package repository

import (
	"context"
	"vk-test-spring/internal/models"
)

type Films interface {
	Create(ctx context.Context, film models.Film) error
	Edit(ctx context.Context, film models.Film) error
	Delete(ctx context.Context, filmId string) error
	GetFilmByName(ctx context.Context, name string) ([]models.Film, error)
	GetFilmByActor(ctx context.Context, actorName string) ([]models.Film, error)
	GetAllFilms(ctx context.Context) ([]models.Film, error)
}

type Actors interface {
	Create(ctx context.Context, actor models.Actor) error
	Edit(ctx context.Context, actor models.Actor) error
	Delete(ctx context.Context, actorId string) error
	GetAllActors(ctx context.Context) ([]models.Actor, error)
}

type Users interface {
	Create(ctx context.Context, user models.User) error
	Delete(ctx context.Context, userId string) error
	Edit(ctx context.Context, user models.User) error
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetAllAdmins(ctx context.Context) ([]models.User, error)
	GetAllCommonUsers(ctx context.Context) ([]models.User, error)
}
