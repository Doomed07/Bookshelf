package bookshelf_service

import (
	"context"
	"fmt"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
	core_errors "github.com/Doomed07/Bookshelf/internal/core/errors"
)

func (s *BookshelfService) GetBooks(
	ctx context.Context,
	userID int,
	read *bool,
	limit, offset *int,
) ([]core_domain.ShelfBookWithBook, error) {
	lim, off, err := core_domain.NormalizePagination(limit, offset)
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

	booksFBS, err := s.bookshelfRepository.GetBooks(ctx, userID, read, lim, off)
	if err != nil {
		return nil, fmt.Errorf("failed to get books from repo: %w", err)
	}

	return booksFBS, nil
}
