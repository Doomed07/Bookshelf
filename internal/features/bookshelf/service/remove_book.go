package bookshelf_service

import (
	"context"
	"fmt"

	core_metrics "github.com/Doomed07/Bookshelf/internal/core/metrics"
)

func (s *BookshelfService) RemoveBook(ctx context.Context, userID, bookID int) error {
	if err := s.bookshelfRepository.RemoveBook(ctx, userID, bookID); err != nil {
		return fmt.Errorf("remove book from repo: %w", err)
	}

	core_metrics.ShelfBooksRemoved.Inc()
	return nil
}
