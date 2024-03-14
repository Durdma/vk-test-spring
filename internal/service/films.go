package service

import (
	"context"
	"vk-test-spring/internal/models"
	"vk-test-spring/internal/repository"
	"vk-test-spring/pkg/logger"
)

type FilmsService struct {
	repo repository.Films
}

func NewFilmsService(repo repository.Films) *FilmsService {
	return &FilmsService{
		repo: repo,
	}
}

func (s *FilmsService) AddNewFilm(ctx context.Context, input FilmInput) error {
	film := models.Film{
		Name:        input.Name,
		Description: input.Description,
		Date:        input.Date,
		Rating:      input.Rating,
	}

	id, err := s.repo.Create(ctx, film)
	if err != nil {
		logger.Error("service 1")
		return err
	}

	if len(input.Actors) > 0 {
		for _, a := range input.Actors {
			err := s.repo.InsertIntoActorFilm(ctx, a, id)
			if err != nil {
				logger.Error("service 2")
				return err
			}
		}
	}

	return nil
}

func (s *FilmsService) EditFilm(ctx context.Context, input FilmInput) error {
	return nil
}
func (s *FilmsService) DeleteFilm(ctx context.Context, name string) error {
	return nil
}
func (s *FilmsService) GetAllFilms(ctx context.Context) ([]models.Film, error) {
	return nil, nil
}
func (s *FilmsService) GetAllFilmsByName(ctx context.Context, name string) ([]models.Film, error) {
	return nil, nil
}
func (s *FilmsService) GetAllFilmsByActor(ctx context.Context, actorsName string) ([]models.Film, error) {
	return nil, nil
}
