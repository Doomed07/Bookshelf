package bookshelf_service

import (
	"context"
	"fmt"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
)

func (s *BookshelfService) GetBook(ctx context.Context, userID, bookID int) (domain.ShelfBookWithBook, error) {
	bfbs, err := s.bookshelfRepository.GetBook(ctx, userID, bookID)
	if err != nil {
		return domain.ShelfBookWithBook{}, fmt.Errorf("get book from repo: %w", err)
	}

	return bfbs, nil
}
