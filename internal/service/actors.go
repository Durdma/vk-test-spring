package service

import (
	"context"
	"github.com/google/uuid"
	"vk-test-spring/internal/models"
	"vk-test-spring/internal/repository"
)

type ActorsService struct {
	repo repository.Actors
}

func NewActorsService(repo repository.Actors) *ActorsService {
	return &ActorsService{
		repo: repo,
	}
}

func (s *ActorsService) AddActor(ctx context.Context, input ActorInput) error {
	actor := models.Actor{
		Name:        input.Name,
		SecondName:  input.SecondName,
		Patronymic:  input.Patronymic,
		Sex:         input.Sex,
		DateOfBirth: input.DateOfBirth,
	}

	id, err := s.repo.Create(ctx, actor)
	if err != nil {
		return err
	}

	if len(input.Films) > 0 {
		for _, f := range input.Films {
			err := s.repo.InsertIntoActorFilm(ctx, id, f)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *ActorsService) resolveActorFilms(actorId uuid.UUID, filmsId []uuid.UUID) error {
	return nil
}

func (s *ActorsService) UpdateActor(ctx context.Context, actor models.Actor) error {
	return nil
}

func (s *ActorsService) DeleteActor(ctx context.Context, actorId string) error {
	return nil
}

func (s *ActorsService) GetAllActors(ctx context.Context) ([]models.Actor, error) {
	return nil, nil
}

func (s *ActorsService) GetActorById(ctx context.Context) (models.Actor, error) {
	return models.Actor{}, nil
}

func (s *ActorsService) GetActorByName(ctx context.Context) ([]models.Actor, error) {
	return nil, nil
}
