package books_transport_http

import (
	"context"
	"net/http"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
	core_http_server "github.com/Doomed07/Bookshelf/internal/core/transport/http/server"
)

type BooksHTTPHandler struct {
	booksService BooksService
}

type BooksService interface {
	GetBooks(ctx context.Context, title, author *string, limit, offset *int) ([]domain.Book, error)
	GetBook(ctx context.Context, id int) (domain.Book, error)
	GetReviews(ctx context.Context, id int, limit, offset *int) ([]domain.Review, error)
}

func NewBooksHTTPHandler(booksService BooksService) *BooksHTTPHandler {
	return &BooksHTTPHandler{
		booksService: booksService,
	}
}

func (h *BooksHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/books",
			Handler: h.GetBooks,
		},
		{
			Method:  http.MethodGet,
			Path:    "/books/{id}",
			Handler: h.GetBook,
		},
		{
			Method:  http.MethodGet,
			Path:    "/books/{id}/reviews",
			Handler: h.GetReviews,
		},
	}
}
