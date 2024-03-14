package httpv1

import (
	"net/http"
	"regexp"
	"vk-test-spring/internal/service"
)

var films = regexp.MustCompile(`^/films/*$`)
var filmsId = regexp.MustCompile(`^/films/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
var filmsWithFilter = regexp.MustCompile(`^/films\?(sort=(name|date|rating)&order=(asc|desc))$`)
var filmsName = regexp.MustCompile(`^/films\?name=.+$`)
var filmsActorName = regexp.MustCompile(`^/films\?actor-name=.+$`)

type FilmsHandler struct {
	filmsService service.Films
}

func NewFilmsHandler(filmsService service.Films) *FilmsHandler {
	return &FilmsHandler{
		filmsService: filmsService,
	}
}

func (h *FilmsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (h *FilmsHandler) AddFilm(w http.ResponseWriter, r *http.Request) {

}

func (h *FilmsHandler) UpdateFilm(w http.ResponseWriter, r *http.Request) {

}

func (h *FilmsHandler) DeleteFilm(w http.ResponseWriter, r *http.Request) {

}

func (h *FilmsHandler) GetAllFilms(w http.ResponseWriter, r *http.Request) {

}

func (h *FilmsHandler) GetFilmsByName(w http.ResponseWriter, r *http.Request) {

}

func (h *FilmsHandler) GetFilmsByActor(w http.ResponseWriter, r *http.Request) {

}
