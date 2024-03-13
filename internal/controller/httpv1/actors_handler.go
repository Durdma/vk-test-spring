package httpv1

import (
	"net/http"
	"vk-test-spring/internal/service"
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
