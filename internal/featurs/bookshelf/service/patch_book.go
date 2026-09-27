package bookshelf_service

import (
	"context"
	"fmt"
	"time"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
)

func (s *BookshelfService) PatchBook(
	ctx context.Context,
	userID, bookID int,
	patch domain.ShelfBookPatch,
) (domain.ShelfBookWithBook, error) {
	shelf, err := s.bookshelfRepository.GetBook(ctx, userID, bookID)
	if err != nil {
		return domain.ShelfBookWithBook{}, fmt.Errorf("get shelfbook from repo: %w", err)
	}

	if err := shelf.ShelfBook.ApplyPatch(patch, time.Now().UTC()); err != nil {
		return domain.ShelfBookWithBook{}, fmt.Errorf("apply shelfBook patch: %w", err)
	}

	if _, err := s.bookshelfRepository.PatchBook(ctx, shelf.ShelfBook); err != nil {
		return domain.ShelfBookWithBook{}, fmt.Errorf("patch shelfBook: %w", err)
	}

	bfsb, err := s.bookshelfRepository.GetBook(ctx, userID, bookID)
	if err != nil {
		return domain.ShelfBookWithBook{}, fmt.Errorf("get shelfbook from repo: %w", err)
	}

	return bfsb, nil
}
