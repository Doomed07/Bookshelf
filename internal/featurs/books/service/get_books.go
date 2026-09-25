package books_service

import (
	"context"
	"fmt"
	"unicode/utf8"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
	core_errors "github.com/Doomed07/Bookshelf/internal/core/errors"
)

func (s *BooksService) GetBooks(
	ctx context.Context,
	title, author *string,
	limit, offset *int,
) ([]domain.Book, error) {
	if title != nil && utf8.RuneCountInString(*title) > 200 {
		return nil, fmt.Errorf("'title' search is longer than 200 symbols: %w", core_errors.ErrInvalidArgument)
	}

	if author != nil && utf8.RuneCountInString(*author) > 100 {
		return nil, fmt.Errorf("'author' search is longer than 100 symbols: %w", core_errors.ErrInvalidArgument)
	}

	lim, off, err := domain.NormalizePagination(limit, offset)
	if err != nil {
		return nil, err
	}

	books, err := s.booksRepository.GetBooks(ctx, title, author, lim, off)
	if err != nil {
		return nil, fmt.Errorf("failed to get books from repo: %w", err)
	}

	return books, nil
}
