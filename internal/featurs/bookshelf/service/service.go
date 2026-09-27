package bookshelf_service

import (
	"context"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
)

type BookshelfService struct {
	bookshelfRepository BookshelfRepository
}

type BookshelfRepository interface {
	AddBook(ctx context.Context, userID, bookID int) (domain.ShelfBook, error)
	GetBook(ctx context.Context, userID, bookID int) (domain.ShelfBookWithBook, error)
	GetBooks(ctx context.Context, userID int, read *bool, limit, offset int) ([]domain.ShelfBookWithBook, error)
	UserExists(ctx context.Context, userID int) (bool, error)
	GetUserActivity(ctx context.Context, userID int, limit, offset int) ([]domain.Event, error)
	PatchBook(ctx context.Context, shelfBook domain.ShelfBook) (domain.ShelfBook, error)
	RemoveBook(ctx context.Context, userID, bookID int) error
}

func NewBookshelfService(bookshelfRepository BookshelfRepository) *BookshelfService {
	return &BookshelfService{
		bookshelfRepository: bookshelfRepository,
	}
}
