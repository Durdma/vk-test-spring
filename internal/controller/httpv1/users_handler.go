package httpv1

import (
	"net/http"
	"vk-test-spring/internal/service"
)

type UsersHandler struct {
	usersService service.Users
}

func NewUsersHandler(usersService service.Users) *UsersHandler {
	return &UsersHandler{
		usersService: usersService,
	}
}

func (h *UsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func (h *UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

}

func (h *UsersHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {

}

func (h *UsersHandler) ChangeRole(w http.ResponseWriter, r *http.Request) {

}
