package users_transport_http

import (
	"context"
	"net/http"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
	core_http_server "github.com/Doomed07/Bookshelf/internal/core/transport/http/server"
)

type UsersHTTPHandler struct {
	usersService UsersService
}

type UsersService interface {
	CreateUser(ctx context.Context, user core_domain.User) (core_domain.User, error)
	GetUsers(ctx context.Context, limit, offset *int) ([]core_domain.User, error)
	GetUser(ctx context.Context, id int) (core_domain.User, error)
	PatchUser(ctx context.Context, id int, patch core_domain.UserPatch) (core_domain.User, error)
	DeleteUser(ctx context.Context, id int) error
}

func NewUsersHTTPHandler(usersService UsersService) *UsersHTTPHandler {
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
			Path:    "/users/{id}",
			Handler: h.GetUser,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/users/{id}",
			Handler: h.PatchUser,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/users/{id}",
			Handler: h.DeleteUser,
		},
	}
}
