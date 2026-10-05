package bookshelf_repository_postgres

import (
	"context"
	"fmt"

	core_errors "github.com/Doomed07/Bookshelf/internal/core/errors"
)

func (r *BookshelfRepository) RemoveBook(ctx context.Context, userID, bookID int) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM bookshelfapp.bookshelf
	WHERE user_id = $1 AND book_id = $2;
	`
	tag, err := r.pool.Exec(ctx, query, userID, bookID)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("book %d on bookshelf of user %d: %w",
			bookID, userID, core_errors.ErrNotFound)
	}

	return nil
}
