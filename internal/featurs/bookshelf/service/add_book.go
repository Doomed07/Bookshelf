package bookshelf_service

import (
	"context"
	"fmt"

	"github.com/Doomed07/Bookshelf/internal/core/core_domain"
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
	return shelfWithBooks, nil
}
