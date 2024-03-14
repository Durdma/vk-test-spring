package httpv1

import (
	"net/http"
	"regexp"
	"vk-test-spring/internal/service"
)

var (
	actors    = regexp.MustCompile(`^/actors/*$`)
	actorId   = regexp.MustCompile(`^/actors/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	actorName = regexp.MustCompile(`^/actors\?name=.+$`)
)

type ActorsHandler struct {
	actorsService service.Actors
}

func NewActorsHandler(actorsService service.Actors) *ActorsHandler {
	return &ActorsHandler{
		actorsService: actorsService,
	}
}

func (h *ActorsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func (h *ActorsHandler) AddActor(w http.ResponseWriter, r *http.Request) {

}

func (h *ActorsHandler) UpdateActor(w http.ResponseWriter, r *http.Request) {

}

func (h *ActorsHandler) DeleteActor(w http.ResponseWriter, r *http.Request) {

}

func (h *ActorsHandler) GetAllActors(w http.ResponseWriter, r *http.Request) {

}
