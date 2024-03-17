package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"regexp"
	"slices"
	"time"
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

type ActorInfo struct {
	Name        string
	SecondName  string
	Patronymic  string
	Sex         string
	DateOfBirth string
}

func (in *ActorInfo) validate() error {
	err := in.validateNamesFields()
	if err != nil {
		return err
	}

	err = in.validateSex()
	if err != nil {
		return err
	}

	err = in.validateDOB()
	if err != nil {
		return err
	}

	return nil
}

func (in *ActorInfo) validateNamesFields() error {
	nameRe := regexp.MustCompile(`^(?:[а-яА-Я]+|[a-zA-Z]+)$`)
	lengthField := 150
	switch {
	case !nameRe.MatchString(in.Name):
		return errors.New(fmt.Sprintf("invalid name format. field must contain"+
			" only characters from the english or russian language, but has: %v", in.Name))
	case !nameRe.MatchString(in.SecondName):
		return errors.New(fmt.Sprintf("invalid second_name format. field must contain"+
			" only characters from the english or russian language, but has: %v", in.SecondName))
	case len(in.Patronymic) > 0 && !nameRe.MatchString(in.Patronymic):
		return errors.New(fmt.Sprintf("invalid patronymic format. field must contain"+
			" only characters from the english or russian language, but has: %v", in.Name))
	case len(in.Name) < 1:
		return errors.New(fmt.Sprintf("input actors's name too short. Length of name must be between 1 and %v,"+
			" but got: %v", lengthField, len(in.Name)))
	case len(in.Name) > lengthField:
		return errors.New(fmt.Sprintf("input actor's name too long. length of name must be between 1 and %v,"+
			" but got: %v", lengthField, len(in.Name)))
	case len(in.SecondName) < 1:
		return errors.New(fmt.Sprintf("input actors's second_name too short. Length of name must be between 1 and %v,"+
			" but got: %v", lengthField, len(in.SecondName)))
	case len(in.SecondName) > lengthField:
		return errors.New(fmt.Sprintf("input actor's second_name too long. length of name must be between 1 and %v,"+
			" but got: %v", lengthField, len(in.SecondName)))
	case len(in.Patronymic) > lengthField:
		return errors.New(fmt.Sprintf("input actor's patronymic too long. length of name must be between 1 and %v,"+
			" but got: %v", lengthField, len(in.Patronymic)))
	default:
		return nil
	}
}

func (in *ActorInfo) validateSex() error {
	if !(in.Sex == "Мужчина" || in.Sex == "Женщина") {
		return errors.New(fmt.Sprintf("invalid value in sex field. field value must be equal to"+
			" 'Мужчина' or 'Женщина', but has: %v", in.Sex))
	}

	return nil
}

func (in *ActorInfo) validateDOB() error {
	d, err := time.Parse(time.DateOnly, in.DateOfBirth)
	if err != nil {
		return err
	}

	before := time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC)

	if d.Before(before) || d.After(time.Now()) {
		return errors.New(fmt.Sprintf("input actors's birthday not in range. date must be in range %v and %v, "+
			"but has: %v", before.Format(time.DateOnly), time.Now().Format(time.DateOnly), in.DateOfBirth))
	}

	return err
}

type ActorCreateInput struct {
	ActorInfo ActorInfo
	Films     []uuid.UUID
}

func (s *ActorsService) AddActor(ctx context.Context, input ActorCreateInput) error {
	err := input.ActorInfo.validate()
	if err != nil {
		return err
	}

	actor := models.Actor{
		Name:        input.ActorInfo.Name,
		SecondName:  input.ActorInfo.SecondName,
		Patronymic:  input.ActorInfo.Patronymic,
		Sex:         input.ActorInfo.Sex,
		DateOfBirth: input.ActorInfo.DateOfBirth,
	}

	err = s.repo.Create(ctx, actor, input.Films)
	if err != nil {
		return err
	}

	return err
}

type ActorUpdateInput struct {
	ActorInfo  ActorInfo
	ID         uuid.UUID
	FilmsToAdd []uuid.UUID
	FilmsToDel []uuid.UUID
}

func (s *ActorsService) UpdateActor(ctx context.Context, input ActorUpdateInput) error {

	actor := models.Actor{
		ID:          input.ID,
		Name:        input.ActorInfo.Name,
		SecondName:  input.ActorInfo.SecondName,
		Patronymic:  input.ActorInfo.Patronymic,
		Sex:         input.ActorInfo.Sex,
		DateOfBirth: input.ActorInfo.DateOfBirth,
	}

	oldActor, err := s.repo.GetActorById(ctx, input.ID)
	if err != nil {
		return err
	}

	actor, err = s.mergeChanges(actor, oldActor)
	if err != nil {
		return err
	}

	actorValidation := ActorInfo{
		Name:        actor.Name,
		SecondName:  actor.SecondName,
		Patronymic:  actor.Patronymic,
		Sex:         actor.Sex,
		DateOfBirth: actor.DateOfBirth,
	}
	err = actorValidation.validate()
	if err != nil {
		return err
	}

	if len(input.FilmsToAdd) > 0 || len(input.FilmsToDel) > 0 {
		err = s.parseFilmsLists(oldActor.Films, input.FilmsToAdd, input.FilmsToDel)
		if err != nil {
			return err
		}
	}

	err = s.repo.Edit(ctx, actor, input.FilmsToAdd, input.FilmsToDel)
	if err != nil {
		return err
	}

	return err
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

func (s *ActorsService) GetActorByName(ctx context.Context, name string) ([]models.Actor, error) {
	return s.repo.GetActorsByName(ctx, name)
}

func (s *ActorsService) mergeChanges(actor models.Actor, oldActor models.Actor) (models.Actor, error) {
	if actor.Name == "" {
		actor.Name = oldActor.Name
	}
	if actor.SecondName == "" {
		actor.SecondName = oldActor.SecondName
	}
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

	for _, f := range filmsToAdd {
		if slices.Contains(currentFilmsUUID, f) {
			return errors.New(fmt.Sprintf("films_to_add film_id that is already in actors_films: %v", f))
		}
	}

	for _, f := range filmsToDel {
		fmt.Println(f)
		fmt.Println(!slices.Contains(currentFilmsUUID, f))
		if !slices.Contains(currentFilmsUUID, f) {
			return errors.New(fmt.Sprintf("films_to_del contains film_id that not in actors_films: %v", f))
		}
	}

	return nil
}
