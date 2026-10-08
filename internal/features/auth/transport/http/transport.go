package auth_transport_http

import (
	"context"
	"net/http"
	"time"

	core_domain "github.com/Doomed07/Bookshelf/internal/core/domain"
	core_http_server "github.com/Doomed07/Bookshelf/internal/core/transport/http/server"
)

type AuthHTTPHandler struct {
	authService AuthService
	ttl         time.Duration
	secure      bool
}

type AuthService interface {
	Register(ctx context.Context, username, email, password string) (core_domain.User, string, error)
	Login(ctx context.Context, login, password string) (core_domain.User, string, error)
	Logout(ctx context.Context, token string) error
}

func NewAuthHTTPHandler(authService AuthService, ttl time.Duration, secure bool) *AuthHTTPHandler {
	return &AuthHTTPHandler{authService: authService, ttl: ttl, secure: secure}
}

func (h *AuthHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{Method: http.MethodPost, Path: "/auth/register", Handler: h.Register},
		{Method: http.MethodPost, Path: "/auth/login", Handler: h.Login},
		{Method: http.MethodPost, Path: "/auth/logout", Handler: h.Logout},
		{Method: http.MethodGet, Path: "/auth/me", Handler: h.Me},
	}
}
