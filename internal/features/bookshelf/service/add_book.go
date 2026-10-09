package bookshelf_service

import (
	"context"
	"fmt"

	core_domain "github.com/Doomed07/Bookshelf/internal/core/domain"
	core_metrics "github.com/Doomed07/Bookshelf/internal/core/metrics"
)

func (s *BookshelfService) AddBook(ctx context.Context, userID, bookID int) (core_domain.ShelfBookWithBook, error) {
	_, err := s.bookshelfRepository.AddBook(ctx, userID, bookID)
	if err != nil {
		return core_domain.ShelfBookWithBook{}, fmt.Errorf("add book from repo: %w", err)
	}

	shelfWithBooks, err := s.bookshelfRepository.GetBook(ctx, userID, bookID)
	if err != nil {
		return core_domain.ShelfBookWithBook{}, fmt.Errorf("get book from repo: %w", err)
	}

	core_metrics.ShelfBooksAdded.Inc()
	return shelfWithBooks, nil
}
