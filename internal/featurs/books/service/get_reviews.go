package books_service

import (
	"context"
	"fmt"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
)

func (s *BooksService) GetReviews(ctx context.Context, id int, limit, offset *int) ([]domain.Review, error) {
	lim, off, err := domain.NormalizePagination(limit, offset)
	if err != nil {
		return nil, err
	}

	if _, err = s.booksRepository.GetBook(ctx, id); err != nil {
		return nil, fmt.Errorf("get book from repo: %w", err)
	}

	reviewsDomain, err := s.booksRepository.GetReviews(ctx, id, lim, off)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews from repo: %w", err)
	}

	return reviewsDomain, nil
}
