package books_service

import (
	"context"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
)

type BooksService struct {
	booksRepository BooksRepository
}

type BooksRepository interface {
	GetBooks(ctx context.Context, title, author *string, limit, offset int) ([]domain.Book, error)
	GetBook(ctx context.Context, id int) (domain.Book, error)
	GetReviews(ctx context.Context, id, limit, offset int) ([]domain.Review, error)
}

func NewBooksService(booksRepository BooksRepository) *BooksService {
	return &BooksService{
		booksRepository: booksRepository,
	}
}
