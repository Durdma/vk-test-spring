package httpv1

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"regexp"
	"strings"
	"vk-test-spring/internal/service"
	"vk-test-spring/pkg/logger"
)

var (
	filmsRe         = regexp.MustCompile(`^/films/*$`)
	filmsIdRe       = regexp.MustCompile(`^/films/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	filmsWithFilter = regexp.MustCompile(`^/films\?(sort=(name|date|rating)&order=(asc|desc))$`)
	//filmsWithFilterV2 ^/films(\?(sort=(name|date|rating)&order=(asc|desc)))?$
	filmsNameRe      = regexp.MustCompile(`^/films\?name=.+$`)
	filmsActorNameRe = regexp.MustCompile(`^/films\?actor-name=.+$`)
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
	case r.Method == http.MethodGet && filmsRe.MatchString(r.URL.Path):
		h.GetAllFilms(w, r)
		return
	case r.Method == http.MethodGet && filmsNameRe.MatchString(r.URL.Path):
		h.GetFilmsByName(w, r)
		return
	case r.Method == http.MethodGet && filmsActorNameRe.MatchString(r.URL.Path):
		h.GetFilmsByActor(w, r)
		return
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

	// TODO добавить приведение даты из строки к дата типу
	err := h.filmsService.AddNewFilm(r.Context(), service.FilmCreateInput{
		Name:        film.Name,
		Description: film.Description,
		Date:        film.Date,
		Rating:      film.Rating,
		Actors:      film.Actors,
	})
	if err != nil {
		logger.Error("Handler")
		// TODO Сделать выбор нужной ошибки и добавить логгирование
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
		ID:          filmId,
		Name:        film.Name,
		Description: film.Description,
		Date:        film.Date,
		Rating:      film.Rating,
		ActorsToAdd: film.ActorsToAdd,
		ActorsToDel: film.ActorsToDel,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *FilmsHandler) getFilmIdFromRequest(r *http.Request) (uuid.UUID, error) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		return uuid.UUID{}, errors.New("error while extracting uuid")
	}

	return uuid.Parse(parts[2])
}

func (h *FilmsHandler) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	filmId, err := h.getFilmIdFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.filmsService.DeleteFilm(r.Context(), filmId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *FilmsHandler) GetAllFilms(w http.ResponseWriter, r *http.Request) {

}

func (h *FilmsHandler) GetFilmsByName(w http.ResponseWriter, r *http.Request) {

}

func (h *FilmsHandler) GetFilmsByActor(w http.ResponseWriter, r *http.Request) {

}
