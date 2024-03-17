package httpv1

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"regexp"
	"strings"
	"vk-test-spring/internal/models"
	"vk-test-spring/internal/service"
)

var (
	filmsRe           = regexp.MustCompile(`^/films/*$`)
	filmsIdRe         = regexp.MustCompile(`^/films/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	filmsWithFilterRe = regexp.MustCompile(`^/films(\?(sort=(name|date|rating)&order=(asc|desc)))?$`)
	//filmsNameRe       = regexp.MustCompile(`^/films\?name=.+$`)
	//filmsActorNameRe  = regexp.MustCompile(`^/films\?actor-name=.+$`)
)

type FilmsHandler struct {
	filmsService service.Films
}

func NewFilmsHandler(filmsService service.Films) *FilmsHandler {
	return &FilmsHandler{
		filmsService: filmsService,
	}
}

func (h *FilmsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && filmsWithFilterRe.MatchString(r.URL.Path):
		params := r.URL.Query()
		switch {
		case params.Get("name") != "":
			h.GetFilmsByName(w, r)
			return
		case params.Get("actor-name") != "":
			h.GetFilmsByActor(w, r)
			return
		default:
			h.GetAllFilms(w, r)
			return
		}
	case r.Method == http.MethodPost && filmsRe.MatchString(r.URL.Path) &&
		r.Context().Value("role").(string) == "администратор":
		h.AddFilm(w, r)
		return
	case r.Method == http.MethodPatch && filmsIdRe.MatchString(r.URL.Path) &&
		r.Context().Value("role").(string) == "администратор":
		h.UpdateFilm(w, r)
		return
	case r.Method == http.MethodDelete && filmsIdRe.MatchString(r.URL.Path) &&
		r.Context().Value("role").(string) == "администратор":
		h.DeleteFilm(w, r)
		return
	default:
		return
	}
}

type FilmCreateInput struct {
	Name        string      `json:"name" binding:"required"`
	Description string      `json:"description" binding:"required"`
	Date        string      `json:"date" binding:"required"`
	Rating      float64     `json:"rating" binding:"required"`
	Actors      []uuid.UUID `json:"actors,omitempty"`
}

func (h *FilmsHandler) AddFilm(w http.ResponseWriter, r *http.Request) {
	var film FilmCreateInput
	if err := json.NewDecoder(r.Body).Decode(&film); err != nil {
		http.Error(w, "error while decoding request body", http.StatusInternalServerError)
		return
	}

	err := h.filmsService.AddNewFilm(r.Context(), service.FilmCreateInput{
		FilmInfo: service.FilmInfo{
			Name:        film.Name,
			Description: film.Description,
			Date:        film.Date,
			Rating:      film.Rating,
		},
		Actors: film.Actors,
	})
	if err != nil {
		switch e := err.(type) {
		case models.CustomError:
			http.Error(w, e.Message, e.Code)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

type FilmUpdateInput struct {
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	Date        string      `json:"date,omitempty"`
	Rating      float64     `json:"rating,omitempty"`
	ActorsToAdd []uuid.UUID `json:"actors_to_add,omitempty"`
	ActorsToDel []uuid.UUID `json:"actors_to_del,omitempty"`
}

func (h *FilmsHandler) UpdateFilm(w http.ResponseWriter, r *http.Request) {
	var film FilmUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&film); err != nil {
		http.Error(w, "error while decoding request body", http.StatusBadRequest)
		return
	}

	filmId, err := h.getFilmIdFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.filmsService.EditFilm(r.Context(), service.FilmUpdateInput{
		ID: filmId,
		FilmInfo: service.FilmInfo{
			Name:        film.Name,
			Description: film.Description,
			Date:        film.Date,
			Rating:      film.Rating,
		},
		ActorsToAdd: film.ActorsToAdd,
		ActorsToDel: film.ActorsToDel,
	})
	if err != nil {
		switch e := err.(type) {
		case models.CustomError:
			http.Error(w, e.Message, e.Code)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (h *FilmsHandler) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	filmId, err := h.getFilmIdFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.filmsService.DeleteFilm(r.Context(), filmId)
	if err != nil {
		switch e := err.(type) {
		case models.CustomError:
			http.Error(w, e.Message, e.Code)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (h *FilmsHandler) GetAllFilms(w http.ResponseWriter, r *http.Request) {
	r = h.getSortParams(r)

	filmsList, err := h.filmsService.GetAllFilms(r.Context())
	if err != nil {
		switch e := err.(type) {
		case models.CustomError:
			http.Error(w, e.Message, e.Code)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	jsonResponse, err := json.Marshal(filmsList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (h *FilmsHandler) GetFilmsByName(w http.ResponseWriter, r *http.Request) {
	var name string
	params := r.URL.Query()
	name = params.Get("name")

	films, err := h.filmsService.GetAllFilmsByName(r.Context(), name)
	if err != nil {
		switch e := err.(type) {
		case models.CustomError:
			http.Error(w, e.Message, e.Code)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	jsonResponse, err := json.Marshal(films)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (h *FilmsHandler) GetFilmsByActor(w http.ResponseWriter, r *http.Request) {
	var actorName string
	params := r.URL.Query()
	actorName = params.Get("name")

	films, err := h.filmsService.GetAllFilmsByActor(r.Context(), actorName)
	if err != nil {
		switch e := err.(type) {
		case models.CustomError:
			http.Error(w, e.Message, e.Code)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	jsonResponse, err := json.Marshal(films)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (h *FilmsHandler) getFilmIdFromRequest(r *http.Request) (uuid.UUID, error) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		return uuid.UUID{}, errors.New("error while extracting uuid")
	}

	return uuid.Parse(parts[2])
}

func (h *FilmsHandler) getSortParams(r *http.Request) *http.Request {
	params := r.URL.Query()
	sortField := params.Get("sort")
	sortOrder := params.Get("order")
	sortOrder = strings.ToUpper(sortOrder)

	if sortField == "" && sortOrder == "" {
		ctx := context.WithValue(r.Context(), "sort", "rating")
		ctx = context.WithValue(ctx, "order", "desc")

		return r.WithContext(ctx)
	}

	ctx := context.WithValue(r.Context(), "sort", sortField)
	ctx = context.WithValue(ctx, "order", sortOrder)

	return r.WithContext(ctx)
}
