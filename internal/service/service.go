package service

import "context"

type FilmInput struct {
}

type Films interface {
	AddNewFilm(ctx context.Context, input FilmInput) error
	EditFilm(ctx context.Context, input FilmInput) error
	DeleteFilm(ctx context.Context, name string) error
}
