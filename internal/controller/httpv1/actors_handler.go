package httpv1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"regexp"
	"strings"
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
	switch {
	case r.Method == http.MethodGet && actors.MatchString(r.URL.Path):
		h.GetAllActors(w, r)
		return
	case r.Method == http.MethodGet && actorId.MatchString(r.URL.Path):
		h.GetActorById(w, r)
		return
	case r.Method == http.MethodPost && actors.MatchString(r.URL.Path) &&
		r.Context().Value("role").(string) == "администратор":
		h.AddActor(w, r)
		return
	case r.Method == http.MethodPatch && actorId.MatchString(r.URL.Path) &&
		r.Context().Value("role").(string) == "администратор":
		h.UpdateActor(w, r)
		return
	case r.Method == http.MethodDelete && actorId.MatchString(r.URL.Path) &&
		r.Context().Value("role").(string) == "администратор":
		h.DeleteActor(w, r)
		return
	default:
		return
	}
}

type ActorCreateInput struct {
	Name        string      `json:"name" binding:"required"`
	SecondName  string      `json:"second_name" binding:"required"`
	Patronymic  string      `json:"patronymic" binding:"required"`
	Sex         string      `json:"sex" binding:"required"`
	DateOfBirth string      `json:"date_of_birth" binding:"required"`
	Films       []uuid.UUID `json:"films,omitempty"`
}

func (h *ActorsHandler) AddActor(w http.ResponseWriter, r *http.Request) {
	var actor ActorCreateInput
	if err := json.NewDecoder(r.Body).Decode(&actor); err != nil {
		http.Error(w, "error while decoding request body", http.StatusBadRequest)
		return
	}

	err := h.actorsService.AddActor(r.Context(), service.ActorCreateInput{
		Name:        actor.Name,
		SecondName:  actor.Name,
		Patronymic:  actor.Patronymic,
		Sex:         actor.Sex,
		DateOfBirth: actor.DateOfBirth,
		Films:       actor.Films,
	})
	if err != nil {
		// TODO Сделать выбор нужной ошибки и добавить логгирование
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type ActorUpdateInput struct {
	Name        string      `json:"name,omitempty"`
	SecondName  string      `json:"second_name,omitempty"`
	Patronymic  string      `json:"patronymic,omitempty"`
	Sex         string      `json:"sex,omitempty"`
	DateOfBirth string      `json:"date_of_birth,omitempty"`
	FilmsToAdd  []uuid.UUID `json:"films_to_add,omitempty"`
	FilmsToDel  []uuid.UUID `json:"films_to_del,omitempty"`
}

func (h *ActorsHandler) UpdateActor(w http.ResponseWriter, r *http.Request) {
	var actor ActorUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&actor); err != nil {
		http.Error(w, "error while decoding request body", http.StatusBadRequest)
		return
	}

	actorId, err := h.getActorIdFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//fmt.Println(actorId)
	//fmt.Println(actor)

	err = h.actorsService.UpdateActor(r.Context(), service.ActorUpdateInput{
		ID:          actorId,
		Name:        actor.Name,
		SecondName:  actor.SecondName,
		Patronymic:  actor.Patronymic,
		Sex:         actor.Sex,
		DateOfBirth: actor.DateOfBirth,
		FilmsToAdd:  actor.FilmsToAdd,
		FilmsToDel:  actor.FilmsToDel,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ActorsHandler) getActorIdFromRequest(r *http.Request) (uuid.UUID, error) {
	parts := strings.Split(r.URL.Path, "/")
	fmt.Println(parts)
	fmt.Println(len(parts))
	if len(parts) != 3 {
		return uuid.UUID{}, errors.New("error while extracting uuid")
	}

	return uuid.Parse(parts[2])
}

func (h *ActorsHandler) DeleteActor(w http.ResponseWriter, r *http.Request) {

}

func (h *ActorsHandler) GetAllActors(w http.ResponseWriter, r *http.Request) {

}

func (h *ActorsHandler) GetActorById(w http.ResponseWriter, r *http.Request) {

}

func (h *ActorsHandler) GetActorByName(w http.ResponseWriter, r *http.Request) {

}
