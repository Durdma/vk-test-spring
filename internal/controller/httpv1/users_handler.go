package httpv1

import (
	"net/http"
	"regexp"
	"vk-test-spring/internal/service"
)

var (
	users   = regexp.MustCompile(`^/users/*$`)
	usersId = regexp.MustCompile(`^/users/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
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

func (h *UsersHandler) GetRole(username string, password string) (string, string, error) {
	return h.usersService.GetUserIdRole(username, password)
}
