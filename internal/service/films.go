package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"slices"
	"time"
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

type FilmInfo struct {
	Name        string
	Description string
	Date        string
	Rating      float64
}

func (in *FilmInfo) validate() error {
	err := in.validateName()
	if err != nil {
		return err
	}

	err = in.validateDescription()
	if err != nil {
		return err
	}

	err = in.validateDate()
	if err != nil {
		return err
	}

	err = in.validateRating()
	if err != nil {
		return err
	}

	return nil
}

func (in *FilmInfo) validateName() error {
	switch {
	case len(in.Name) > 150:
		return errors.New(fmt.Sprintf("input film's name too long. length of name must be between 1 and 150,"+
			" but got: %v", len(in.Name)))
	case len(in.Name) == 0:
		return errors.New("input film's name is empty")
	default:
		return nil
	}
}

func (in *FilmInfo) validateDescription() error {
	if len(in.Description) > 1000 {
		return errors.New(fmt.Sprintf("input film's description too long. length of name must be between 1 and 1000,"+
			" but got: %v", len(in.Description)))
	}
	if len(in.Description) < 1 {
		return errors.New("empty film's description")
	}

	return nil
}

func (in *FilmInfo) validateDate() error {
	d, err := time.Parse(time.DateOnly, in.Date)
	if err != nil {
		return err
	}

	before := time.Date(1895, time.December, 28, 0, 0, 0, 0, time.UTC)

	if d.Before(before) {
		return errors.New(fmt.Sprintf("input film's date not in range. date cant be earlier %v, "+
			"but has: %v", before.Format(time.DateOnly), in.Date))
	}

	return err
}

func (in *FilmInfo) validateRating() error {
	switch {
	case in.Rating < 0:
		return errors.New(fmt.Sprintf("input films's rating is negative. rating value must be in range between 0 and 10,"+
			" but got: %v", in.Rating))
	case in.Rating > 10:
		return errors.New(fmt.Sprintf("input films's rating is too big. rating value must be in range between 0 and 10,"+
			" but got: %v", in.Rating))
	default:
		return nil
	}
}

type FilmCreateInput struct {
	FilmInfo FilmInfo
	Actors   []uuid.UUID
}

func (s *FilmsService) AddNewFilm(ctx context.Context, input FilmCreateInput) error {
	err := input.FilmInfo.validate()
	if err != nil {
		return models.CustomError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	film := models.Film{
		Name:        input.FilmInfo.Name,
		Description: input.FilmInfo.Description,
		Date:        input.FilmInfo.Date,
		Rating:      input.FilmInfo.Rating,
	}

	err = s.repo.Create(ctx, film, input.Actors)
	if err != nil {
		return err
	}

	return err
}

type FilmUpdateInput struct {
	ID          uuid.UUID
	FilmInfo    FilmInfo
	ActorsToAdd []uuid.UUID
	ActorsToDel []uuid.UUID
}

func (s *FilmsService) EditFilm(ctx context.Context, input FilmUpdateInput) error {
	film := models.Film{
		ID:          input.ID,
		Name:        input.FilmInfo.Name,
		Description: input.FilmInfo.Description,
		Date:        input.FilmInfo.Date,
		Rating:      input.FilmInfo.Rating,
	}

	oldFilm, err := s.repo.GetFilmById(ctx, input.ID)
	if err != nil {
		return err
	}

	film, err = s.mergeChanges(film, oldFilm)
	if err != nil {
		return err
	}

	filmValidation := FilmInfo{
		Name:        film.Name,
		Description: film.Description,
		Date:        film.Date,
		Rating:      film.Rating,
	}
	err = filmValidation.validate()
	if err != nil {
		return models.CustomError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	if len(input.ActorsToAdd) > 0 || len(input.ActorsToDel) > 0 {
		err = s.parseActorsLists(oldFilm.Actors, input.ActorsToAdd, input.ActorsToDel)
		if err != nil {
			return models.CustomError{Code: http.StatusBadRequest, Message: err.Error()}
		}
	}

	err = s.repo.Update(ctx, film, input.ActorsToAdd, input.ActorsToDel)
	if err != nil {
		return err
	}

	return err
}
func (s *FilmsService) DeleteFilm(ctx context.Context, filmId uuid.UUID) error {
	return s.repo.Delete(ctx, filmId)
}
func (s *FilmsService) GetAllFilms(ctx context.Context) ([]models.Film, error) {
	return s.repo.GetAllFilms(ctx)
}
func (s *FilmsService) GetAllFilmsByName(ctx context.Context, name string) ([]models.Film, error) {
	return s.repo.GetFilmByName(ctx, name)
}
func (s *FilmsService) GetAllFilmsByActor(ctx context.Context, actorsName string) ([]models.Film, error) {

	return s.repo.GetFilmByActor(ctx, actorsName)
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
		if !slices.Contains(currentActorsUUID, a) {
			return errors.New(fmt.Sprintf("actors_to_del contains actor_id that not in film_actors: %v", a))
		}
	}

	return nil
}
