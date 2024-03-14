package controller

import (
	"net/http"
	"vk-test-spring/internal/controller/httpv1"
	"vk-test-spring/internal/service"
)

type Handler struct {
	filmsHandler  FilmsHandler
	actorsHandler ActorsHandler
	usersHandler  UsersHandler
}

type ActorsHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	AddActor(w http.ResponseWriter, r *http.Request)
	UpdateActor(w http.ResponseWriter, r *http.Request)
	DeleteActor(w http.ResponseWriter, r *http.Request)
	GetAllActors(w http.ResponseWriter, r *http.Request)
}

type FilmsHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	AddFilm(w http.ResponseWriter, r *http.Request)
	UpdateFilm(w http.ResponseWriter, r *http.Request)
	DeleteFilm(w http.ResponseWriter, r *http.Request)
	GetAllFilms(w http.ResponseWriter, r *http.Request)
	GetFilmsByName(w http.ResponseWriter, r *http.Request)
	GetFilmsByActor(w http.ResponseWriter, r *http.Request)
}

type UsersHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	ChangeRole(w http.ResponseWriter, r *http.Request)
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Init(services *service.Services) *http.ServeMux {
	mux := http.NewServeMux()

	h.filmsHandler = httpv1.NewFilmsHandler(services.Films)
	h.actorsHandler = httpv1.NewActorsHandler(services.Actors)
	h.usersHandler = httpv1.NewUsersHandler(services.Users)

	h.initAPI(mux)

	return mux
}

func (h *Handler) initAPI(router *http.ServeMux) {
	router.Handle("/films", h.filmsHandler)
	router.Handle("/actors", h.actorsHandler)
	router.Handle("/users", h.usersHandler)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
