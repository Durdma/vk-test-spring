package httpv1

import (
	"net/http"
	"vk-test-spring/internal/service"
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
