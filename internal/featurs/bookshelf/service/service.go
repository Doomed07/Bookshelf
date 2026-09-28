package bookshelf_service

import (
	"context"

	"github.com/Doomed07/Bookshelf/internal/core/core_domain"
)

type BookshelfService struct {
	bookshelfRepository BookshelfRepository
}

type BookshelfRepository interface {
	AddBook(ctx context.Context, userID, bookID int) (core_domain.ShelfBook, error)
	GetBook(ctx context.Context, userID, bookID int) (core_domain.ShelfBookWithBook, error)
	GetBooks(ctx context.Context, userID int, read *bool, limit, offset int) ([]core_domain.ShelfBookWithBook, error)
	UserExists(ctx context.Context, userID int) (bool, error)
	GetUserActivity(ctx context.Context, userID int, limit, offset int) ([]core_domain.Event, error)
	PatchBook(ctx context.Context, shelfBook core_domain.ShelfBook) (core_domain.ShelfBook, error)
	RemoveBook(ctx context.Context, userID, bookID int) error
}

func NewBookshelfService(bookshelfRepository BookshelfRepository) *BookshelfService {
	return &BookshelfService{
		bookshelfRepository: bookshelfRepository,
	}
}
