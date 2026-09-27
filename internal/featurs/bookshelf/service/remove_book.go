package bookshelf_service

import (
	"context"
	"fmt"
)

func (s *BookshelfService) RemoveBook(ctx context.Context, userID, bookID int) error {
	if err := s.bookshelfRepository.RemoveBook(ctx, userID, bookID); err != nil {
		return fmt.Errorf("remove book from repo: %w", err)
	}
	return nil
}
