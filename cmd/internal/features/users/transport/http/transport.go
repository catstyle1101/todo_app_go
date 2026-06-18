package users_transport_http

import (
	"context"
	"net/http"

	"github.com/catstyle1101/todo_app_go/cmd/internal/core/domain"
	core_http_server "github.com/catstyle1101/todo_app_go/cmd/internal/core/transport/http/server"
)

type UsersHTTPHandler struct {
	usersService UsersService
}

type UsersService interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)

	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.User, error)

	GetUser(
		ctx context.Context,
		userID int,
	) (domain.User, error)

	DeleteUser(
		ctx context.Context,
		userID int,
	) error

	PatchUser(
		ctx context.Context,
		id int,
		patch domain.UserPatch,

	) (domain.User, error)
}

func NewUsersHttpHandler(usersService UsersService) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: usersService,
	}
}

func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: h.GetUsers,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/{user_id}",
			Handler: h.GetUser,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/users/{user_id}",
			Handler: h.DeleteUser,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/users/{user_id}",
			Handler: h.PatchUser,
		},
	}
}
