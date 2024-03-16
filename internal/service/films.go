package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"slices"
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

func (s *FilmsService) AddNewFilm(ctx context.Context, input FilmCreateInput) error {
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
		err = s.addFilmActors(ctx, id, input.Actors)
	}

	return err
}

func (s *FilmsService) addFilmActors(ctx context.Context, filmId uuid.UUID, actorsId []uuid.UUID) error {
	for _, a := range actorsId {
		err := s.repo.InsertIntoActorFilm(ctx, a, filmId)
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO WE ARE HERE
func (s *FilmsService) EditFilm(ctx context.Context, input FilmUpdateInput) error {
	film := models.Film{
		ID:          input.ID,
		Name:        input.Name,
		Description: input.Description,
		Date:        input.Date,
		Rating:      input.Rating,
	}

	oldFilm, err := s.repo.GetFilmById(ctx, input.ID)
	if err != nil {
		return err
	}

	film, err = s.mergeChanges(film, oldFilm)
	if err != nil {
		return err
	}

	fmt.Println("===")
	fmt.Println(film.ID)
	fmt.Println(oldFilm.ID)
	fmt.Println("===")

	if len(input.ActorsToAdd) > 0 || len(input.ActorsToDel) > 0 {
		err = s.parseActorsLists(oldFilm.Actors, input.ActorsToAdd, input.ActorsToDel)
		if err != nil {
			return err
		}
	}

	err = s.repo.Update(ctx, film, input.ActorsToAdd, input.ActorsToDel)
	if err != nil {
		return err
	}

	return err
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

func (s *FilmsService) mergeChanges(film models.Film, oldFilm models.Film) (models.Film, error) {
	if film.Name == "" {
		film.Name = oldFilm.Name
	}

	if film.Description == "" {
		film.Description = oldFilm.Description
	}

	if film.Rating == 0 {
		film.Rating = oldFilm.Rating
	}

	if film.Date == "" {
		film.Date = oldFilm.Date
	}

	return film, nil
}

// TODO Rewrite this function to parse only toDel and toAdd parse, DB will return err, if no id in list
func (s *FilmsService) parseActorsLists(currentActors []models.FilmActors, actorsToAdd []uuid.UUID, actorsToDel []uuid.UUID) error {
	for _, a := range actorsToAdd {
		if slices.Contains(actorsToDel, a) {
			return errors.New(fmt.Sprintf("actors_to_add and actors_to_del contains same actor_id: %v", a))
		}
	}

	currentActorsUUID := make([]uuid.UUID, 0, len(currentActors))
	for _, a := range currentActors {
		currentActorsUUID = append(currentActorsUUID, a.ID)
	}

	for _, a := range actorsToAdd {
		if slices.Contains(currentActorsUUID, a) {
			return errors.New(fmt.Sprintf("actors_to_add contains actor_id that is already in film_actors: %v", a))
		}
	}

	for _, a := range actorsToDel {
		if slices.Contains(currentActorsUUID, a) {
			return errors.New(fmt.Sprintf("actors_to_del contains actor_id that not in film_actors: %v", a))
		}
	}

	return nil
}

func (s *FilmsService) removeFilmActors(ctx context.Context, filmId uuid.UUID, actorsId []uuid.UUID) error {
	for _, a := range actorsId {
		err := s.repo.DeleteFromActorFilm(ctx, filmId, a)
		if err != nil {
			return err
		}
	}

	return nil
}
