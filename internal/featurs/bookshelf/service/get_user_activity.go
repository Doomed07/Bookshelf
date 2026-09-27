package bookshelf_service

import (
	"context"
	"fmt"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
	core_errors "github.com/Doomed07/Bookshelf/internal/core/errors"
)

func (s *BookshelfService) GetUserActivity(ctx context.Context, userID int, limit, offset *int) ([]domain.Event, error) {
	lim, off, err := domain.NormalizePagination(limit, offset)
	if err != nil {
		return nil, err
	}

	exists, err := s.bookshelfRepository.UserExists(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("check user exists: %w", err)
	}

	if !exists {
		return nil, fmt.Errorf("user %d: %w", userID, core_errors.ErrNotFound)
	}

	events, err := s.bookshelfRepository.GetUserActivity(ctx, userID, lim, off)
	if err != nil {
		return nil, fmt.Errorf("get events from repo: %w", err)
	}

	return events, nil
}
