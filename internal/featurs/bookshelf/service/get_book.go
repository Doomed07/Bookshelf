package bookshelf_service

import (
	"context"
	"fmt"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
)

func (s *BookshelfService) GetBook(ctx context.Context, userID, bookID int) (core_domain.ShelfBookWithBook, error) {
	bfbs, err := s.bookshelfRepository.GetBook(ctx, userID, bookID)
	if err != nil {
		return core_domain.ShelfBookWithBook{}, fmt.Errorf("get book from repo: %w", err)
	}

	return bfbs, nil
}
