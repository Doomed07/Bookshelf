package bookshelf_transport_http

import (
	"context"
	"net/http"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
	core_http_server "github.com/Doomed07/Bookshelf/internal/core/transport/http/server"
)

type BookshelfHTTPHandler struct {
	bookshelfService BookshelfService
}

type BookshelfService interface {
	AddBook(ctx context.Context, userID, bookID int) (domain.ShelfBookWithBook, error)
	GetBook(ctx context.Context, userID, bookID int) (domain.ShelfBookWithBook, error)
	GetBooks(ctx context.Context, userID int, read *bool, limit, offset *int) ([]domain.ShelfBookWithBook, error)
	GetUserActivity(ctx context.Context, userID int, limit, offset *int) ([]domain.Event, error)
	PatchBook(ctx context.Context, userID, bookID int, patch domain.ShelfBookPatch) (domain.ShelfBookWithBook, error)
	RemoveBook(ctx context.Context, userID, bookID int) error
}

func NewBookshelfHTTPHandler(bookshelfService BookshelfService) *BookshelfHTTPHandler {
	return &BookshelfHTTPHandler{
		bookshelfService: bookshelfService,
	}
}

func (h *BookshelfHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users/{user_id}/bookshelf",
			Handler: h.AddBook,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/{user_id}/bookshelf/{book_id}",
			Handler: h.GetBook,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/{user_id}/bookshelf",
			Handler: h.GetBooks,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/{user_id}/activity",
			Handler: h.GetUserActivity,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/users/{user_id}/bookshelf/{book_id}",
			Handler: h.PatchBook,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/users/{user_id}/bookshelf/{book_id}",
			Handler: h.RemoveBook,
		},
	}
}
