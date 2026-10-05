package books_service

import (
	"context"
	"fmt"

	core_domain "github.com/Doomed07/Bookshelf/internal/core/domain"
)

func (s *BooksService) GetRecentReviews(ctx context.Context, limit, offset *int) ([]core_domain.ReviewWithBook, error) {
	lim, off, err := core_domain.NormalizePagination(limit, offset)
	if err != nil {
		return nil, err
	}

	reviews, err := s.booksRepository.GetRecentReviews(ctx, lim, off)
	if err != nil {
		return nil, fmt.Errorf("get recent reviews from repo: %w", err)
	}

	return reviews, nil
}
