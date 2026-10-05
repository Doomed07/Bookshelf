package books_service

import (
	"context"

	core_domain "github.com/Doomed07/Bookshelf/internal/core/domain"
)

type BooksService struct {
	booksRepository BooksRepository
}

type BooksRepository interface {
	GetBooks(ctx context.Context, title, author *string, limit, offset int) ([]core_domain.Book, error)
	GetBook(ctx context.Context, id int) (core_domain.Book, error)
	GetReviews(ctx context.Context, id, limit, offset int) ([]core_domain.Review, error)
	GetRecentReviews(ctx context.Context, limit, offset int) ([]core_domain.ReviewWithBook, error)
}

func NewBooksService(booksRepository BooksRepository) *BooksService {
	return &BooksService{
		booksRepository: booksRepository,
	}
}
