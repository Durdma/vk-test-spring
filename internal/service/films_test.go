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

type MockFilmRepository struct {
	mock.Mock
}

//type Films interface {
//	Create(ctx context.Context, film models.Film, actors []uuid.UUID) error
//	Update(ctx context.Context, film models.Film, actorsToAdd []uuid.UUID, actorsToDel []uuid.UUID) error
//	Delete(ctx context.Context, filmId uuid.UUID) error
//	GetFilmByName(ctx context.Context, name string) ([]models.Film, error)
//	GetFilmByActor(ctx context.Context, actorName string) ([]models.Film, error)
//	GetAllFilms(ctx context.Context) ([]models.Film, error)
//	GetFilmById(ctx context.Context, filmId uuid.UUID) (models.Film, error)
//}

func (m *MockFilmRepository) Create(ctx context.Context, film models.Film, actors []uuid.UUID) error {
	args := m.Called(ctx, film, actors)
	return args.Error(0)
}

func (m *MockFilmRepository) Update(ctx context.Context, film models.Film, actorsToAdd []uuid.UUID, actorsToDel []uuid.UUID) error {
	args := m.Called(ctx, film, actorsToAdd, actorsToDel)
	return args.Error(0)
}

func (m *MockFilmRepository) Delete(ctx context.Context, filmId uuid.UUID) error {
	args := m.Called(ctx, filmId)
	return args.Error(0)
}

func (m *MockFilmRepository) GetFilmByName(ctx context.Context, name string) ([]models.Film, error) {
	args := m.Called(ctx, name)
	return args.Get(0).([]models.Film), args.Error(0)
}

func (m *MockFilmRepository) GetFilmByActor(ctx context.Context, actorName string) ([]models.Film, error) {
	args := m.Called(ctx, actorName)
	return args.Get(0).([]models.Film), args.Error(0)
}

func (m *MockFilmRepository) GetAllFilms(ctx context.Context) ([]models.Film, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Film), args.Error(0)
}

func (m *MockFilmRepository) GetFilmById(ctx context.Context, filmId uuid.UUID) (models.Film, error) {
	args := m.Called(ctx, filmId)
	return args.Get(0).(models.Film), args.Error(1)
}

func TestFilmsService_AddNewFilm(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := new(MockFilmRepository)
		filmService := FilmsService{repo: repo}

		input := FilmCreateInput{
			FilmInfo: FilmInfo{
				Name:        "Film",
				Description: "films description",
				Date:        "2000-01-01",
				Rating:      5.2,
			},
			Actors: nil,
		}

		repo.On("Create", context.Background(), models.Film{
			Name:        input.FilmInfo.Name,
			Description: input.FilmInfo.Description,
			Date:        input.FilmInfo.Date,
			Rating:      input.FilmInfo.Rating,
		}, input.Actors).Return(nil)

		err := filmService.AddNewFilm(context.Background(), input)

		assert.NoError(t, err)

		repo.AssertCalled(t, "Create", context.Background(), models.Film{
			Name:        input.FilmInfo.Name,
			Description: input.FilmInfo.Description,
			Date:        input.FilmInfo.Date,
			Rating:      input.FilmInfo.Rating,
		}, input.Actors)
	})

	t.Run("empty film name", func(t *testing.T) {
		repo := new(MockFilmRepository)
		filmService := FilmsService{repo: repo}

		input := FilmCreateInput{
			FilmInfo: FilmInfo{
				Name:        "",
				Description: "films description",
				Date:        "2000-01-01",
				Rating:      5.2,
			},
			Actors: nil,
		}

		repo.On("Create", mock.Anything, mock.Anything).Return(nil)

		err := filmService.AddNewFilm(context.Background(), input)

		assert.Error(t, err)
		assert.EqualError(t, err, "input film's name is empty")
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("too long name", func(t *testing.T) {
		repo := new(MockFilmRepository)
		filmService := FilmsService{repo: repo}

		input := FilmCreateInput{
			FilmInfo: FilmInfo{
				Name:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				Description: "films description",
				Date:        "2000-01-01",
				Rating:      5.2,
			},
			Actors: nil,
		}

		repo.On("Create", mock.Anything, mock.Anything).Return(nil)

		err := filmService.AddNewFilm(context.Background(), input)

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("input film's name too long. length of name must be between 1 and 150,"+
			" but got: %v", len(input.FilmInfo.Name)))
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("too long description", func(t *testing.T) {
		repo := new(MockFilmRepository)
		filmService := FilmsService{repo: repo}

		input := FilmCreateInput{
			FilmInfo: FilmInfo{
				Name: "test film description",
				Description: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
					"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
					"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
					"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
					"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
					"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
					"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
					"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
					"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
					"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				Date:   "2000-01-01",
				Rating: 5.2,
			},
			Actors: nil,
		}

		repo.On("Create", mock.Anything, mock.Anything).Return(nil)

		err := filmService.AddNewFilm(context.Background(), input)

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("input film's description too long. length of name must be between 1 and 1000,"+
			" but got: %v", len(input.FilmInfo.Description)))
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("empty description", func(t *testing.T) {
		repo := new(MockFilmRepository)
		filmService := FilmsService{repo: repo}

		input := FilmCreateInput{
			FilmInfo: FilmInfo{
				Name:        "test film description",
				Description: "",
				Date:        "2000-01-01",
				Rating:      5.2,
			},
			Actors: nil,
		}

		repo.On("Create", mock.Anything, mock.Anything).Return(nil)

		err := filmService.AddNewFilm(context.Background(), input)

		assert.Error(t, err)
		assert.EqualError(t, err, "empty film's description")
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("incorrect date", func(t *testing.T) {
		repo := new(MockFilmRepository)
		filmService := FilmsService{repo: repo}

		input := FilmCreateInput{
			FilmInfo: FilmInfo{
				Name:        "test film description",
				Description: "description",
				Date:        "41-31-2011",
				Rating:      5.2,
			},
			Actors: nil,
		}

		repo.On("Create", mock.Anything, mock.Anything).Return(nil)

		err := filmService.AddNewFilm(context.Background(), input)

		assert.Error(t, err)

		repo.AssertNotCalled(t, "Create")
	})

	t.Run("less than 0 rating", func(t *testing.T) {
		repo := new(MockFilmRepository)
		filmService := FilmsService{repo: repo}

		input := FilmCreateInput{
			FilmInfo: FilmInfo{
				Name:        "test film description",
				Description: "description",
				Date:        "2000-01-01",
				Rating:      -4.3,
			},
			Actors: nil,
		}

		repo.On("Create", mock.Anything, mock.Anything).Return(nil)

		err := filmService.AddNewFilm(context.Background(), input)

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("input films's rating is negative. rating value must be in range between 0 and 10,"+
			" but got: %v", input.FilmInfo.Rating))
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("more than 10 rating", func(t *testing.T) {
		repo := new(MockFilmRepository)
		filmService := FilmsService{repo: repo}

		input := FilmCreateInput{
			FilmInfo: FilmInfo{
				Name:        "test film description",
				Description: "description",
				Date:        "2000-01-01",
				Rating:      10.1,
			},
			Actors: nil,
		}

		repo.On("Create", mock.Anything, mock.Anything).Return(nil)

		err := filmService.AddNewFilm(context.Background(), input)

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("input films's rating is too big. rating value must be in range between 0 and 10,"+
			" but got: %v", input.FilmInfo.Rating))
		repo.AssertNotCalled(t, "Create")
	})
}

func TestFilmsService_EditFilm(t *testing.T) {
	t.Run("Success Update", func(t *testing.T) {
		repo := new(MockFilmRepository)
		filmService := FilmsService{repo: repo}

		input := FilmUpdateInput{
			FilmInfo: FilmInfo{
				Name:        "test film",
				Description: "description",
				Date:        "2000-01-01",
				Rating:      5.5,
			},
			ID:          uuid.UUID{},
			ActorsToAdd: nil,
			ActorsToDel: nil,
		}

		expectedFromGetFilmById := models.Film{
			ID:          input.ID,
			Name:        "test film",
			Description: "description",
			Date:        "2000-01-01",
			Rating:      5.5,
			Actors:      nil,
		}

		repo.On("GetFilmById", context.Background(), input.ID).Return(expectedFromGetFilmById, nil)

		repo.On("Update", context.Background(), expectedFromGetFilmById, input.ActorsToAdd, input.ActorsToDel).Return(nil)

		err := filmService.EditFilm(context.Background(), input)

		assert.NoError(t, err)

		repo.AssertCalled(t, "GetFilmById", context.Background(), input.ID)

		repo.AssertCalled(t, "Update", context.Background(), expectedFromGetFilmById, input.ActorsToAdd, input.ActorsToDel)
	})

	t.Run("record not found", func(t *testing.T) {
		repo := new(MockFilmRepository)
		filmService := FilmsService{repo: repo}

		repo.On("GetFilmById", context.Background(), uuid.Nil).Return(models.Film{}, errors.New("record not found"))

		err := filmService.EditFilm(context.Background(), FilmUpdateInput{ID: uuid.Nil})

		assert.Equal(t, err.Error(), "record not found")

		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("Good Parse List to add", func(t *testing.T) {
		repo := new(MockFilmRepository)
		filmService := FilmsService{repo: repo}

		input := FilmUpdateInput{
			FilmInfo: FilmInfo{
				Name:        "test film",
				Description: "description",
				Date:        "2000-01-01",
				Rating:      5.5,
			},
			ID:          uuid.UUID{},
			ActorsToAdd: []uuid.UUID{uuid.New()},
			ActorsToDel: nil,
		}

		expectedFromGetFilmById := models.Film{
			ID:          input.ID,
			Name:        "test film",
			Description: "description",
			Date:        "2000-01-01",
			Rating:      5.5,
			Actors:      nil,
		}

		repo.On("GetFilmById", context.Background(), input.ID).Return(expectedFromGetFilmById, nil)

		repo.On("Update", context.Background(), expectedFromGetFilmById, input.ActorsToAdd, input.ActorsToDel).Return(nil)

		err := filmService.EditFilm(context.Background(), input)

		assert.NoError(t, err)

		repo.AssertCalled(t, "Update", context.Background(), expectedFromGetFilmById, input.ActorsToAdd, input.ActorsToDel)
	})

	t.Run("good parse list to del", func(t *testing.T) {
		repo := new(MockFilmRepository)
		filmService := FilmsService{repo: repo}

		todel := uuid.New()

		input := FilmUpdateInput{
			FilmInfo: FilmInfo{
				Name:        "test film",
				Description: "description",
				Date:        "2000-01-01",
				Rating:      5.5,
			},
			ID:          uuid.UUID{},
			ActorsToAdd: nil,
			ActorsToDel: []uuid.UUID{todel},
		}

		expectedFromGetFilmById := models.Film{
			ID:          input.ID,
			Name:        "test film",
			Description: "description",
			Date:        "2000-01-01",
			Rating:      5.5,
			Actors: []models.FilmActors{{
				ID:   todel,
				Name: "test",
			}},
		}

		repo.On("GetFilmById", context.Background(), input.ID).Return(expectedFromGetFilmById, nil)

		expectedEdited := expectedFromGetFilmById
		expectedEdited.Actors = nil

		repo.On("Update", context.Background(), expectedEdited, input.ActorsToAdd, input.ActorsToDel).Return(nil)

		err := filmService.EditFilm(context.Background(), input)

		assert.NoError(t, err)

		repo.AssertCalled(t, "GetFilmById", context.Background(), input.ID)

		repo.AssertCalled(t, "Update", context.Background(), expectedEdited, input.ActorsToAdd, input.ActorsToDel)
	})

	t.Run("Bad parse list to add", func(t *testing.T) {

	})
}
