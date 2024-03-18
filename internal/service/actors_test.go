package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"vk-test-spring/internal/models"
)

type MockActorRepository struct {
	mock.Mock
}

func (m *MockActorRepository) Create(ctx context.Context, actor models.Actor, actorFilms []uuid.UUID) error {
	args := m.Called(ctx, actor, actorFilms)
	return args.Error(0)
}

func (m *MockActorRepository) Edit(ctx context.Context, actor models.Actor, filmsToAdd []uuid.UUID, filmsToDel []uuid.UUID) error {
	args := m.Called(ctx, actor, filmsToAdd, filmsToDel)
	return args.Error(0)
}

func (m *MockActorRepository) Delete(ctx context.Context, actorId uuid.UUID) error {
	args := m.Called(ctx, actorId)
	return args.Error(0)
}

func (m *MockActorRepository) GetAllActors(ctx context.Context) ([]models.Actor, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Actor), args.Error(1)
}

func (m *MockActorRepository) GetActorsByName(ctx context.Context, name string) ([]models.Actor, error) {
	args := m.Called(ctx, name)
	return args.Get(0).([]models.Actor), args.Error(1)
}

func (m *MockActorRepository) GetActorById(ctx context.Context, actorId uuid.UUID) (models.Actor, error) {
	args := m.Called(ctx, actorId)
	return args.Get(0).(models.Actor), args.Error(1)
}

// TODO rewrite test cases like 1st
func TestActorService_AddActor(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		input := ActorCreateInput{
			ActorInfo: ActorInfo{
				Name:        "Maxim",
				SecondName:  "Sizask",
				Patronymic:  "Edu",
				Sex:         "Мужчина",
				DateOfBirth: "2000-01-01",
			},
			Films: nil,
		}

		repo.On("Create", context.Background(), models.Actor{
			Name:        input.ActorInfo.Name,
			SecondName:  input.ActorInfo.SecondName,
			Patronymic:  input.ActorInfo.Patronymic,
			Sex:         input.ActorInfo.Sex,
			DateOfBirth: input.ActorInfo.DateOfBirth,
		}, input.Films).Return(nil)

		err := actorService.AddActor(context.Background(), input)

		assert.NoError(t, err)

		repo.AssertCalled(t, "Create", context.Background(), models.Actor{
			Name:        input.ActorInfo.Name,
			SecondName:  input.ActorInfo.SecondName,
			Patronymic:  input.ActorInfo.Patronymic,
			Sex:         input.ActorInfo.Sex,
			DateOfBirth: input.ActorInfo.DateOfBirth,
		}, input.Films)
	})

	t.Run("Not valid Name", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		input := ActorCreateInput{ActorInfo: ActorInfo{
			Name:        "МаксимS",
			SecondName:  "Sizask",
			Patronymic:  "Edu",
			Sex:         "Мужчина",
			DateOfBirth: "2000-01-01",
		},
			Films: nil,
		}

		repo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := actorService.AddActor(context.Background(), input)

		assert.Error(t, err)
		assert.EqualError(t, err, "invalid name format. field must contain"+
			" only characters from the english or russian language, but has: МаксимS")
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("Not valid SecondName", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		input := ActorCreateInput{ActorInfo: ActorInfo{
			Name:        "Максим",
			SecondName:  "SizaskЫ",
			Patronymic:  "Edu",
			Sex:         "Мужчина",
			DateOfBirth: "2000-01-01",
		},
			Films: nil,
		}

		repo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := actorService.AddActor(context.Background(), input)

		assert.Error(t, err)
		assert.EqualError(t, err, "invalid second_name format. field must contain"+
			" only characters from the english or russian language, but has: SizaskЫ")
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("Not valid Patronymic", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		input := ActorCreateInput{ActorInfo: ActorInfo{
			Name:        "Максим",
			SecondName:  "Sizask",
			Patronymic:  "Eduард",
			Sex:         "Мужчина",
			DateOfBirth: "2000-01-01",
		},
			Films: nil,
		}

		repo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := actorService.AddActor(context.Background(), input)

		assert.Error(t, err)
		assert.EqualError(t, err, "invalid patronymic format. field must contain"+
			" only characters from the english or russian language, but has: Eduард")
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("Empty Name", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		input := ActorCreateInput{ActorInfo: ActorInfo{
			Name:        "",
			SecondName:  "Sizask",
			Patronymic:  "Edu",
			Sex:         "Мужчина",
			DateOfBirth: "2000-01-01",
		},
			Films: nil,
		}

		repo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := actorService.AddActor(context.Background(), input)

		assert.Error(t, err)
		assert.EqualError(t, err, "input actors's name too short. Length of name must be between 1 and 150,"+
			" but got: 0")
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("Empty SecondName", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		input := ActorCreateInput{ActorInfo: ActorInfo{
			Name:        "Максим",
			SecondName:  "",
			Patronymic:  "Edu",
			Sex:         "Мужчина",
			DateOfBirth: "2000-01-01",
		},
			Films: nil,
		}

		repo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := actorService.AddActor(context.Background(), input)

		assert.Error(t, err)
		assert.EqualError(t, err, "input actors's second_name too short. Length of name must be between 1 and 150,"+
			" but got: 0")
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("Empty Patronymic", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		input := ActorCreateInput{ActorInfo: ActorInfo{
			Name:        "Maxim",
			SecondName:  "Sizask",
			Patronymic:  "",
			Sex:         "Мужчина",
			DateOfBirth: "2000-01-01",
		},
			Films: nil,
		}

		repo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := actorService.AddActor(context.Background(), input)

		assert.NoError(t, err)
		repo.AssertCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("Long Name", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		input := ActorCreateInput{ActorInfo: ActorInfo{
			Name:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			SecondName:  "Sizask",
			Patronymic:  "Edu",
			Sex:         "Мужчина",
			DateOfBirth: "2000-01-01",
		},
			Films: nil,
		}

		repo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := actorService.AddActor(context.Background(), input)

		assert.Error(t, err)
		assert.EqualError(t, err, "input actor's name too long. length of name must be between 1 and 150,"+
			" but got: 174")
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("Long SecondName", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		input := ActorCreateInput{ActorInfo: ActorInfo{
			Name:        "Max",
			SecondName:  "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			Patronymic:  "Edu",
			Sex:         "Мужчина",
			DateOfBirth: "2000-01-01",
		},
			Films: nil,
		}

		repo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := actorService.AddActor(context.Background(), input)

		assert.Error(t, err)
		assert.EqualError(t, err, "input actor's second_name too long. length of name must be between 1 and 150,"+
			" but got: 174")
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("Long Patronymic", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		input := ActorCreateInput{ActorInfo: ActorInfo{
			Name:        "Max",
			SecondName:  "Sizask",
			Patronymic:  "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			Sex:         "Мужчина",
			DateOfBirth: "2000-01-01",
		},
			Films: nil,
		}

		repo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := actorService.AddActor(context.Background(), input)

		assert.Error(t, err)
		assert.EqualError(t, err, "input actor's patronymic too long. length of name must be between 1 and 150,"+
			" but got: 174")
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("Correct Sex", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		input := ActorCreateInput{ActorInfo: ActorInfo{
			Name:        "Максим",
			SecondName:  "Siz",
			Patronymic:  "Edu",
			Sex:         "Мужчина",
			DateOfBirth: "2000-01-01",
		},
			Films: nil,
		}

		repo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := actorService.AddActor(context.Background(), input)

		assert.NoError(t, err)
		repo.AssertCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("InCorrect Sex", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		input := ActorCreateInput{ActorInfo: ActorInfo{
			Name:        "Максим",
			SecondName:  "Siz",
			Patronymic:  "Edu",
			Sex:         "Мужчинаasd",
			DateOfBirth: "2000-01-01",
		},
			Films: nil,
		}

		repo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := actorService.AddActor(context.Background(), input)

		assert.Error(t, err)
		assert.EqualError(t, err, "invalid value in sex field. field value must be equal to"+
			" 'Мужчина' or 'Женщина', but has: Мужчинаasd")
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("Correct DOB", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		input := ActorCreateInput{ActorInfo: ActorInfo{
			Name:        "Максим",
			SecondName:  "Siz",
			Patronymic:  "Edu",
			Sex:         "Мужчина",
			DateOfBirth: "2000-01-01",
		},
			Films: nil,
		}

		repo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := actorService.AddActor(context.Background(), input)

		assert.NoError(t, err)
		repo.AssertCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("InCorrect DOB", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		input := ActorCreateInput{ActorInfo: ActorInfo{
			Name:        "Максим",
			SecondName:  "Siz",
			Patronymic:  "Edu",
			Sex:         "Мужчина",
			DateOfBirth: "41-31-2011",
		},
			Films: nil,
		}

		repo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := actorService.AddActor(context.Background(), input)

		assert.Error(t, err)
		//assert.EqualError(t, err, "invalid value in sex field. field value must be equal to"+
		//	" 'Мужчина' or 'Женщина', but has: Мужчинаasd")
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("Not in range DOB", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		input := ActorCreateInput{ActorInfo: ActorInfo{
			Name:        "Максим",
			SecondName:  "Siz",
			Patronymic:  "Edu",
			Sex:         "Мужчина",
			DateOfBirth: "1885-01-01",
		},
			Films: nil,
		}

		repo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := actorService.AddActor(context.Background(), input)

		assert.Error(t, err)
		assert.EqualError(t, err, "input actors's birthday not in range. date must be in range 1900-01-01 and 2024-03-18, "+
			"but has: 1885-01-01")
		repo.AssertNotCalled(t, "Create")
	})
}

func TestActorsService_UpdateActor(t *testing.T) {
	t.Run("Success Update 2", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		input := ActorUpdateInput{
			ActorInfo: ActorInfo{
				Name:        "Maxim",
				SecondName:  "Sizask",
				Patronymic:  "Edu",
				Sex:         "Мужчина",
				DateOfBirth: "2000-01-12",
			},
			ID:         uuid.UUID{},
			FilmsToAdd: nil,
			FilmsToDel: nil,
		}

		expectedFromGetActorById := models.Actor{
			ID:          input.ID,
			Name:        "Maxim",
			SecondName:  "Sizask",
			Patronymic:  "Edu",
			Sex:         "Мужчина",
			DateOfBirth: "2000-01-12",
			Films:       nil,
		}

		// Устанавливаем ожидание для вызова метода GetActorById
		repo.On("GetActorById", context.Background(), input.ID).Return(expectedFromGetActorById, nil)

		// Устанавливаем ожидание для вызова метода Edit
		repo.On("Edit", context.Background(), expectedFromGetActorById, input.FilmsToAdd, input.FilmsToDel).Return(nil)

		// Вызываем метод UpdateActor
		err := actorService.UpdateActor(context.Background(), input)

		// Проверяем, что нет ошибок
		assert.NoError(t, err)

		// Проверяем, что метод GetActorById был вызван с правильными аргументами
		repo.AssertCalled(t, "GetActorById", context.Background(), input.ID)

		// Проверяем, что метод Edit был вызван с правильными аргументами
		repo.AssertCalled(t, "Edit", context.Background(), expectedFromGetActorById, input.FilmsToAdd, input.FilmsToDel)
	})

	t.Run("Record Not Found", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		repo.On("GetActorById", context.Background(), uuid.Nil).Return(models.Actor{ID: uuid.Nil}, errors.New("record not found"))

		// Вызываем метод UpdateActor
		err := actorService.UpdateActor(context.Background(), ActorUpdateInput{ID: uuid.Nil})
		// Проверяем, что ошибка возвращается и она соответствует ожидаемой ошибке "запись не найдена".
		//assert.Error(t, err)
		assert.Equal(t, "record not found", err.Error())

		// Проверяем, что метод Edit не вызывается, потому что запись не найдена.
		repo.AssertNotCalled(t, "Edit", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("Good Parse List To add", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		input := ActorUpdateInput{
			ActorInfo: ActorInfo{
				Name:        "Maxim",
				SecondName:  "Sizask",
				Patronymic:  "Edu",
				Sex:         "Мужчина",
				DateOfBirth: "2000-01-12",
			},
			ID:         uuid.UUID{},
			FilmsToAdd: []uuid.UUID{uuid.New()},
			FilmsToDel: nil,
		}

		expectedGetActor := models.Actor{
			ID:          input.ID,
			Name:        "Maxim",
			SecondName:  "Sizask",
			Patronymic:  "Edu",
			Sex:         "Мужчина",
			DateOfBirth: "2000-01-12",
			Films:       nil,
		}

		repo.On("GetActorById", context.Background(), input.ID).Return(expectedGetActor, nil)
		repo.On("Edit", context.Background(), expectedGetActor, input.FilmsToAdd, input.FilmsToDel).Return(nil)

		err := actorService.UpdateActor(context.Background(), input)

		assert.NoError(t, err)

		repo.AssertCalled(t, "Edit", context.Background(), expectedGetActor, input.FilmsToAdd, input.FilmsToDel)
	})

	t.Run("Good Parse List To del", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		todel := uuid.New()

		input := ActorUpdateInput{
			ActorInfo: ActorInfo{
				Name:        "Maxim",
				SecondName:  "Sizask",
				Patronymic:  "Edu",
				Sex:         "Мужчина",
				DateOfBirth: "2000-01-12",
			},
			ID:         uuid.UUID{},
			FilmsToAdd: nil,
			FilmsToDel: []uuid.UUID{todel},
		}

		expectedGetActor := models.Actor{
			ID:          input.ID,
			Name:        "Maxim",
			SecondName:  "Sizask",
			Patronymic:  "Edu",
			Sex:         "Мужчина",
			DateOfBirth: "2000-01-12",
			Films: []models.ActorFilm{{
				ID:   todel,
				Name: "test",
			}},
		}

		// Устанавливаем ожидание для вызова метода GetActorById
		repo.On("GetActorById", context.Background(), input.ID).Return(expectedGetActor, nil)

		// Создаём копию актёра с пустым списком фильмов
		expectedEditedActor := expectedGetActor
		expectedEditedActor.Films = nil

		// Устанавливаем ожидание для вызова метода Edit
		repo.On("Edit", context.Background(), expectedEditedActor, input.FilmsToAdd, input.FilmsToDel).Return(nil)

		// Вызываем метод UpdateActor
		err := actorService.UpdateActor(context.Background(), input)

		// Проверяем, что нет ошибок
		assert.NoError(t, err)

		// Проверяем, что метод GetActorById был вызван с правильными аргументами
		repo.AssertCalled(t, "GetActorById", context.Background(), input.ID)

		// Проверяем, что метод Edit был вызван с правильными аргументами
		repo.AssertCalled(t, "Edit", context.Background(), expectedEditedActor, input.FilmsToAdd, input.FilmsToDel)
	})

	t.Run("Bad Parse List To add", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		added := uuid.New()

		input := ActorUpdateInput{
			ActorInfo: ActorInfo{
				Name:        "Maxim",
				SecondName:  "Sizask",
				Patronymic:  "Edu",
				Sex:         "Мужчина",
				DateOfBirth: "2000-01-12",
			},
			ID:         uuid.UUID{},
			FilmsToAdd: []uuid.UUID{added},
			FilmsToDel: nil,
		}

		expectedGetActor := models.Actor{
			ID:          input.ID,
			Name:        "Maxim",
			SecondName:  "Sizask",
			Patronymic:  "Edu",
			Sex:         "Мужчина",
			DateOfBirth: "2000-01-12",
			Films: []models.ActorFilm{{
				ID:   added,
				Name: "test",
			}},
		}

		repo.On("GetActorById", context.Background(), input.ID).Return(expectedGetActor, nil)
		//repo.On("Edit", context.Background(), expectedGetActor, input.FilmsToAdd, input.FilmsToDel).Return(nil)

		err := actorService.UpdateActor(context.Background(), input)

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("films_to_add film_id that is already in actors_films: %v", added))

		repo.AssertNotCalled(t, "Edit", context.Background(), expectedGetActor, input.FilmsToAdd, input.FilmsToDel)
	})

	t.Run("Bad Parse List To del", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		deleted := uuid.New()

		input := ActorUpdateInput{
			ActorInfo: ActorInfo{
				Name:        "Maxim",
				SecondName:  "Sizask",
				Patronymic:  "Edu",
				Sex:         "Мужчина",
				DateOfBirth: "2000-01-12",
			},
			ID:         uuid.UUID{},
			FilmsToAdd: nil,
			FilmsToDel: []uuid.UUID{deleted},
		}

		expectedGetActor := models.Actor{
			ID:          input.ID,
			Name:        "Maxim",
			SecondName:  "Sizask",
			Patronymic:  "Edu",
			Sex:         "Мужчина",
			DateOfBirth: "2000-01-12",
			Films:       nil,
		}

		repo.On("GetActorById", context.Background(), input.ID).Return(expectedGetActor, nil)
		//repo.On("Edit", context.Background(), expectedGetActor, input.FilmsToAdd, input.FilmsToDel).Return(nil)

		err := actorService.UpdateActor(context.Background(), input)

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("films_to_del contains film_id that not in actors_films: %v", deleted))

		repo.AssertNotCalled(t, "Edit", context.Background(), expectedGetActor, input.FilmsToAdd, input.FilmsToDel)
	})

	t.Run("Bad Parse List To del and To add", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		deleted := uuid.New()

		input := ActorUpdateInput{
			ActorInfo: ActorInfo{
				Name:        "Maxim",
				SecondName:  "Sizask",
				Patronymic:  "Edu",
				Sex:         "Мужчина",
				DateOfBirth: "2000-01-12",
			},
			ID:         uuid.UUID{},
			FilmsToAdd: []uuid.UUID{deleted},
			FilmsToDel: []uuid.UUID{deleted},
		}

		expectedGetActor := models.Actor{
			ID:          input.ID,
			Name:        "Maxim",
			SecondName:  "Sizask",
			Patronymic:  "Edu",
			Sex:         "Мужчина",
			DateOfBirth: "2000-01-12",
			Films:       nil,
		}

		repo.On("GetActorById", context.Background(), input.ID).Return(expectedGetActor, nil)
		//repo.On("Edit", context.Background(), expectedGetActor, input.FilmsToAdd, input.FilmsToDel).Return(nil)

		err := actorService.UpdateActor(context.Background(), input)

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("films_to_add and films_to_del contains same film_id: %v", deleted))

		repo.AssertNotCalled(t, "Edit", context.Background(), expectedGetActor, input.FilmsToAdd, input.FilmsToDel)
	})

	t.Run("Merge Changes", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		//deleted := uuid.New()

		input := ActorUpdateInput{
			ActorInfo: ActorInfo{
				Name:        "",
				SecondName:  "",
				Patronymic:  "",
				Sex:         "",
				DateOfBirth: "",
			},
			ID:         uuid.UUID{},
			FilmsToAdd: nil,
			FilmsToDel: nil,
		}

		expectedGetActor := models.Actor{
			ID:          input.ID,
			Name:        "Maxim",
			SecondName:  "Sizask",
			Patronymic:  "Edu",
			Sex:         "Мужчина",
			DateOfBirth: "2000-01-12",
			Films:       nil,
		}

		repo.On("GetActorById", context.Background(), input.ID).Return(expectedGetActor, nil)
		repo.On("Edit", context.Background(), expectedGetActor, input.FilmsToAdd, input.FilmsToDel).Return(nil)

		err := actorService.UpdateActor(context.Background(), input)

		assert.NoError(t, err)
		//assert.EqualError(t, err, fmt.Sprintf("films_to_add and films_to_del contains same film_id: %v", deleted))

		repo.AssertCalled(t, "Edit", context.Background(), expectedGetActor, input.FilmsToAdd, input.FilmsToDel)
	})
}

func TestActorsService_DeleteActor(t *testing.T) {
	t.Run("Success Del", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		toDel := uuid.New()
		// Устанавливаем ожидание для вызова метода GetActorById
		repo.On("Delete", context.Background(), toDel).Return(nil)

		// Вызываем метод UpdateActor
		err := actorService.DeleteActor(context.Background(), toDel)

		// Проверяем, что нет ошибок
		assert.NoError(t, err)

		repo.AssertCalled(t, "Delete", context.Background(), toDel)
	})

	t.Run("Record Not Found", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		toDel := uuid.New()
		// Устанавливаем ожидание для вызова метода GetActorById
		repo.On("Delete", context.Background(), toDel).Return(errors.New("record not found"))

		// Вызываем метод UpdateActor
		err := actorService.DeleteActor(context.Background(), toDel)

		// Проверяем, что нет ошибок
		assert.Error(t, err)
		assert.EqualError(t, err, "record not found")
		repo.AssertCalled(t, "Delete", context.Background(), toDel)
	})
}

func TestActorsService_GetAllActors(t *testing.T) {
	t.Run("Success get all", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		// Устанавливаем ожидание для вызова метода GetActorById
		repo.On("GetAllActors", context.Background()).Return([]models.Actor{{ID: uuid.New()}}, nil)

		// Вызываем метод UpdateActor
		_, err := actorService.GetAllActors(context.Background())

		// Проверяем, что нет ошибок
		assert.NoError(t, err)

		repo.AssertCalled(t, "GetAllActors", context.Background())
	})

	t.Run("Records Not Found", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		// Устанавливаем ожидание для вызова метода GetActorById
		repo.On("GetAllActors", context.Background()).Return([]models.Actor{}, errors.New("records not found"))

		// Вызываем метод UpdateActor
		_, err := actorService.GetAllActors(context.Background())

		// Проверяем, что нет ошибок
		assert.Error(t, err)
		assert.EqualError(t, err, "records not found")

		repo.AssertCalled(t, "GetAllActors", context.Background())
	})
}

func TestActorsService_GetActorById(t *testing.T) {
	t.Run("Success get actor by Id", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		id := uuid.New()

		// Устанавливаем ожидание для вызова метода GetActorById
		repo.On("GetActorById", context.Background(), id).Return(models.Actor{ID: id}, nil)

		// Вызываем метод UpdateActor
		_, err := actorService.GetActorById(context.Background(), id)

		// Проверяем, что нет ошибок
		assert.NoError(t, err)

		repo.AssertCalled(t, "GetActorById", context.Background(), id)
	})

	t.Run("not found actor by id", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		id := uuid.New()

		// Устанавливаем ожидание для вызова метода GetActorById
		repo.On("GetActorById", context.Background(), id).Return(models.Actor{}, errors.New("record not found"))

		// Вызываем метод UpdateActor
		_, err := actorService.GetActorById(context.Background(), id)

		// Проверяем, что нет ошибок
		assert.Error(t, err)
		assert.EqualError(t, err, "record not found")

		repo.AssertCalled(t, "GetActorById", context.Background(), id)
	})
}

func TestActorsService_GetActorByName(t *testing.T) {
	t.Run("Success get actor by name", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		name := "test"

		// Устанавливаем ожидание для вызова метода GetActorById
		repo.On("GetActorsByName", context.Background(), name).Return([]models.Actor{{Name: name}}, nil)

		// Вызываем метод UpdateActor
		_, err := actorService.GetActorByName(context.Background(), name)

		// Проверяем, что нет ошибок
		assert.NoError(t, err)

		repo.AssertCalled(t, "GetActorsByName", context.Background(), name)
	})

	t.Run("not found actor by id", func(t *testing.T) {
		repo := new(MockActorRepository)
		actorService := ActorsService{repo: repo}

		name := "test"

		// Устанавливаем ожидание для вызова метода GetActorById
		repo.On("GetActorsByName", context.Background(), name).Return([]models.Actor{}, errors.New("records not found"))

		// Вызываем метод UpdateActor
		_, err := actorService.GetActorByName(context.Background(), name)

		// Проверяем, что нет ошибок
		assert.Error(t, err)
		assert.EqualError(t, err, "records not found")

		repo.AssertCalled(t, "GetActorsByName", context.Background(), name)
	})
}
