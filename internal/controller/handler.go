package controller

import (
	"context"
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
	GetActorById(w http.ResponseWriter, r *http.Request)
	GetActorByName(w http.ResponseWriter, r *http.Request)
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
	GetRole(username string, password string) (string, string, error)
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
	router.Handle("/films", h.usersAuth(h.filmsHandler)) // добавить это на handlers ниже
	router.Handle("/films/", h.usersAuth(h.filmsHandler))
	router.Handle("/actors", h.usersAuth(h.actorsHandler))
	router.Handle("/actors/", h.usersAuth(h.actorsHandler))
	router.Handle("/users", h.usersHandler)
	router.Handle("/users/", h.usersHandler)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) usersAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userId, role, err := h.usersHandler.GetRole(username, password)
		if err != nil {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userId)
		ctx = context.WithValue(r.Context(), "role", role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
