package bookshelf_service

import (
	"context"
	"fmt"

	"github.com/Doomed07/Bookshelf/internal/core/core_domain"
)

func (s *BookshelfService) PatchBook(
	ctx context.Context,
	userID, bookID int,
	patch core_domain.ShelfBookPatch,
) (core_domain.ShelfBookWithBook, error) {
	shelf, err := s.bookshelfRepository.GetBook(ctx, userID, bookID)
	if err != nil {
		return core_domain.ShelfBookWithBook{}, fmt.Errorf("get shelfbook from repo: %w", err)
	}

	if err := shelf.ShelfBook.ApplyPatch(patch); err != nil {
		return core_domain.ShelfBookWithBook{}, fmt.Errorf("apply shelfBook patch: %w", err)
	}

	if _, err := s.bookshelfRepository.PatchBook(ctx, shelf.ShelfBook); err != nil {
		return core_domain.ShelfBookWithBook{}, fmt.Errorf("patch shelfBook: %w", err)
	}

	bfsb, err := s.bookshelfRepository.GetBook(ctx, userID, bookID)
	if err != nil {
		return core_domain.ShelfBookWithBook{}, fmt.Errorf("get shelfbook from repo: %w", err)
	}

	return bfsb, nil
}
