package service

import (
	"context"
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
