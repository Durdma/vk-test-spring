package service

import (
	"context"
	"vk-test-spring/internal/models"
	"vk-test-spring/internal/repository"
)

type FilmsService struct {
	repo repository.Films
}

func NewFilmsService(repo repository.Films) *FilmsService {
	return &FilmsService{
		repo: repo,
	}
}

func (f *FilmsService) AddNewFilm(ctx context.Context, input FilmInput) error {
	return nil
}

func (f *FilmsService) EditFilm(ctx context.Context, input FilmInput) error {
	return nil
}
func (f *FilmsService) DeleteFilm(ctx context.Context, name string) error {
	return nil
}
func (f *FilmsService) GetAllFilms(ctx context.Context) ([]models.Film, error) {
	return nil, nil
}
func (f *FilmsService) GetAllFilmsByName(ctx context.Context, name string) ([]models.Film, error) {
	return nil, nil
}
func (f *FilmsService) GetAllFilmsByActor(ctx context.Context, actorsName string) ([]models.Film, error) {
	return nil, nil
}
