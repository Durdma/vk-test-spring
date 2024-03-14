package httpv1

import (
	"net/http"
	"regexp"
	"vk-test-spring/internal/service"
)

var (
	films           = regexp.MustCompile(`^/films/*$`)
	filmsId         = regexp.MustCompile(`^/films/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	filmsWithFilter = regexp.MustCompile(`^/films\?(sort=(name|date|rating)&order=(asc|desc))$`)
	filmsName       = regexp.MustCompile(`^/films\?name=.+$`)
	filmsActorName  = regexp.MustCompile(`^/films\?actor-name=.+$`)
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
	case r.Method == http.MethodGet && films.MatchString(r.URL.Path):
		h.GetAllFilms(w, r)
		return
	case r.Method == http.MethodGet && filmsName.MatchString(r.URL.Path):
		h.GetFilmsByName(w, r)
		return
	case r.Method == http.MethodGet && filmsActorName.MatchString(r.URL.Path):
		h.GetFilmsByActor(w, r)
		return
	case r.Method == http.MethodPost && films.MatchString(r.URL.Path) &&
		r.Context().Value("role").(string) == "администратор":
		h.AddFilm(w, r)
		return
	case r.Method == http.MethodPatch && filmsId.MatchString(r.URL.Path) &&
		r.Context().Value("role").(string) == "администратор":
		h.UpdateFilm(w, r)
		return
	case r.Method == http.MethodDelete && filmsId.MatchString(r.URL.Path) &&
		r.Context().Value("role").(string) == "администратор":
		h.DeleteFilm(w, r)
		return
	default:
		return
	}
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
