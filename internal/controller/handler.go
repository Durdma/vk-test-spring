package controller

import (
	"context"
	"github.com/rs/zerolog"
	_ "github.com/swaggo/files"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"time"
	"vk-test-spring/internal/controller/httpv1"
	"vk-test-spring/internal/service"
	"vk-test-spring/pkg/logger"
)

type Handler struct {
	filmsHandler  FilmsHandler
	actorsHandler ActorsHandler
	usersHandler  UsersHandler
	logger        zerolog.Logger
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

func (h *Handler) Init(services *service.Services, logs zerolog.Logger) *http.ServeMux {
	mux := http.NewServeMux()

	h.filmsHandler = httpv1.NewFilmsHandler(services.Films)
	h.actorsHandler = httpv1.NewActorsHandler(services.Actors)
	h.usersHandler = httpv1.NewUsersHandler(services.Users)

	h.logger = logs

	h.initAPI(mux)

	return mux
}

func (h *Handler) initAPI(router *http.ServeMux) {
	router.Handle("/films", h.logs(h.usersAuth(h.filmsHandler))) // добавить это на handlers ниже
	router.Handle("/films/", h.logs(h.usersAuth(h.filmsHandler)))
	router.Handle("/actors", h.logs(h.usersAuth(h.actorsHandler)))
	router.Handle("/actors/", h.logs(h.usersAuth(h.actorsHandler)))
	router.Handle("/users", h.logs(h.usersHandler))
	router.Handle("/users/", h.logs(h.usersHandler))
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
		ctx = context.WithValue(ctx, "role", role)

		l := ctx.Value("logger").(*zerolog.Event)
		l.Str("user_id", userId).Str("user_name", username).Str("role", role)

		ctx = context.WithValue(ctx, "logger", l)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) logs(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := logger.NewWrapResponseWriter(w, r.ProtoMajor)
		rec := httptest.NewRecorder()

		ctx := r.Context()

		path := r.URL.EscapedPath()

		reqData, _ := httputil.DumpRequest(r, true)

		logg := h.logger.Log().Timestamp().Str("path", path).Bytes("request_data", reqData)

		defer func(begin time.Time) {
			status := ww.Status()

			tookMs := time.Since(begin).Milliseconds()
			logg.Int64("took", tookMs).Int("status_code", status).Msgf("[%d] %s http request for %s took %dms",
				status, r.Method, path, tookMs)
		}(time.Now())

		ctx = context.WithValue(ctx, "logger", logg)
		next.ServeHTTP(rec, r.WithContext(ctx))

		for k, v := range rec.Header() {
			ww.Header()[k] = v
		}

		ww.WriteHeader(rec.Code)
		rec.Body.WriteTo(ww)
	})
}
