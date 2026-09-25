package books_service

import (
	"context"
	"fmt"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
)

func (s *BooksService) GetBook(ctx context.Context, id int) (domain.Book, error) {
	bookDomain, err := s.booksRepository.GetBook(ctx, id)
	if err != nil {
		return domain.Book{}, fmt.Errorf("get book from repo: %w", err)
	}

	return bookDomain, nil
}
