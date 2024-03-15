package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"slices"
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

func (s *ActorsService) AddActor(ctx context.Context, input ActorCreateInput) error {
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
		err = s.addActorFilms(ctx, id, input.Films)
	}

	return err
}

func (s *ActorsService) UpdateActor(ctx context.Context, input ActorUpdateInput) error {
	actor := models.Actor{
		ID:          input.ID,
		Name:        input.Name,
		SecondName:  input.SecondName,
		Patronymic:  input.Patronymic,
		Sex:         input.Sex,
		DateOfBirth: input.DateOfBirth,
	}

	oldActor, err := s.repo.GetActorById(ctx, input.ID)
	if err != nil {
		return err
	}

	// TODO Validate input
	actor, err = s.mergeChanges(actor, oldActor)
	if err != nil {
		return err
	}

	if len(input.FilmsToAdd) > 0 || len(input.FilmsToDel) > 0 {
		err = s.parseFilmsLists(actor.Films, input.FilmsToAdd, input.FilmsToDel)
		if err != nil {
			return err
		}
	}

	err = s.repo.Edit(ctx, actor)
	if err != nil {
		return err
	}

	if len(input.FilmsToAdd) > 0 {
		err = s.addActorFilms(ctx, actor.ID, input.FilmsToAdd)
		if err != nil {
			return err
		}
	}

	if len(input.FilmsToDel) > 0 {
		err = s.removeActorFilms(ctx, actor.ID, input.FilmsToDel)
		if err != nil {
			return err
		}
	}

	return err
}

func (s *ActorsService) mergeChanges(actor models.Actor, oldActor models.Actor) (models.Actor, error) {
	if actor.Name == "" {
		actor.Name = oldActor.Name
	}
	if actor.SecondName == "" {
		actor.SecondName = oldActor.SecondName
	}
	// TODO Что делать если у актера больше нет отчества???
	if actor.Patronymic == "" {
		actor.Patronymic = oldActor.Patronymic
	}
	if actor.Sex == "" {
		actor.Sex = oldActor.Sex
	}
	if actor.DateOfBirth == "" {
		actor.DateOfBirth = oldActor.DateOfBirth
	}

	return actor, nil
}

func (s *ActorsService) parseFilmsLists(currentFilms []models.ActorFilm, filmsToAdd []uuid.UUID, filmsToDel []uuid.UUID) error {
	for _, f := range filmsToAdd {
		if slices.Contains(filmsToDel, f) {
			return errors.New(fmt.Sprintf("films_to_add and films_to_del contains same film_id: %v", f))
		}
	}

	currentFilmsUUID := make([]uuid.UUID, 0, len(currentFilms))

	for _, f := range currentFilms {
		currentFilmsUUID = append(currentFilmsUUID, f.ID)
	}

	fmt.Println(currentFilmsUUID)

	for _, f := range filmsToAdd {
		if slices.Contains(currentFilmsUUID, f) {
			return errors.New(fmt.Sprintf("films_to_add film_id that is already in actors_films: %v", f))
		}
	}

	for _, f := range filmsToDel {
		if slices.Contains(currentFilmsUUID, f) {
			fmt.Println(slices.Contains(currentFilmsUUID, f))
			return errors.New(fmt.Sprintf("films_to_del contains film_id that not in actors_films: %v", f))
		}
	}

	return nil
}

func (s *ActorsService) DeleteActor(ctx context.Context, actorId uuid.UUID) error {
	return s.repo.Delete(ctx, actorId)
}

func (s *ActorsService) GetAllActors(ctx context.Context) ([]models.Actor, error) {
	return s.repo.GetAllActors(ctx)
}

func (s *ActorsService) GetActorById(ctx context.Context, actorId uuid.UUID) (models.Actor, error) {

	return s.repo.GetActorById(ctx, actorId)
}

func (s *ActorsService) GetActorByName(ctx context.Context) ([]models.Actor, error) {
	return nil, nil
}

func (s *ActorsService) addActorFilms(ctx context.Context, actorId uuid.UUID, filmsId []uuid.UUID) error {
	for _, f := range filmsId {
		err := s.repo.InsertIntoActorFilm(ctx, actorId, f)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *ActorsService) removeActorFilms(ctx context.Context, actorId uuid.UUID, filmsId []uuid.UUID) error {
	for _, f := range filmsId {
		err := s.repo.DeleteFromActorFilm(ctx, actorId, f)
		if err != nil {
			return err
		}
	}

	return nil
}
