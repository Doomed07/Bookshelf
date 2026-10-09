package bookshelf_service

import (
	"context"
	"fmt"

	core_domain "github.com/Doomed07/Bookshelf/internal/core/domain"
	core_metrics "github.com/Doomed07/Bookshelf/internal/core/metrics"
)

func (s *BookshelfService) PatchBook(
	ctx context.Context,
	userID, bookID int,
	patch core_domain.ShelfBookPatch,
) (core_domain.ShelfBookWithBook, error) {
	shelf, err := s.bookshelfRepository.GetBook(ctx, userID, bookID)
	if err != nil {
		return core_domain.ShelfBookWithBook{},
			fmt.Errorf("get shelfbook from repo: %w", err)
	}

	wasRead := shelf.ShelfBook.Read
	hadReview := hasText(shelf.ShelfBook.Review)

	if err := shelf.ShelfBook.ApplyPatch(patch); err != nil {
		return core_domain.ShelfBookWithBook{},
			fmt.Errorf("apply shelfBook patch: %w", err)
	}

	if _, err := s.bookshelfRepository.PatchBook(ctx, shelf.ShelfBook); err != nil {
		return core_domain.ShelfBookWithBook{},
			fmt.Errorf("patch shelfBook: %w", err)
	}

	if !wasRead && shelf.ShelfBook.Read {
		core_metrics.BooksMarkedRead.Inc()
	}
	if !hadReview && hasText(shelf.ShelfBook.Review) {
		core_metrics.ReviewsPublished.Inc()
	}

	bfsb, err := s.bookshelfRepository.GetBook(ctx, userID, bookID)
	if err != nil {
		return core_domain.ShelfBookWithBook{},
			fmt.Errorf("get shelfbook from repo: %w", err)
	}

	return bfsb, nil
}

func hasText(s *string) bool {
	return s != nil && *s != ""
}
