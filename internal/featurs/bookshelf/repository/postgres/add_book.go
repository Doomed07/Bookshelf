package bookshelf_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Doomed07/Bookshelf/internal/core/core_domain"
	core_errors "github.com/Doomed07/Bookshelf/internal/core/errors"
	core_postgres_pool "github.com/Doomed07/Bookshelf/internal/core/repository/postgres/pool"
)

func (r *BookshelfRepository) AddBook(ctx context.Context, userID, bookID int) (core_domain.ShelfBook, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO bookshelfapp.bookshelf (user_id, book_id)
	VALUES ($1, $2)
	RETURNING user_id, book_id, version, read, rating, review, added_at, read_at;
	`

	row := r.pool.QueryRow(ctx, query, userID, bookID)

	var shelf ShelfBookModel
	err := row.Scan(
		&shelf.UserID,
		&shelf.BookID,
		&shelf.Version,
		&shelf.Read,
		&shelf.Rating,
		&shelf.Review,
		&shelf.AddedAt,
		&shelf.ReadAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, core_postgres_pool.ErrUniqueViolation):
			return core_domain.ShelfBook{}, fmt.Errorf(
				"book %d already on bookshelf of user %d: %w", bookID, userID, core_errors.ErrConflict)
		case errors.Is(err, core_postgres_pool.ErrForeignKeyViolation):
			return core_domain.ShelfBook{}, fmt.Errorf(
				"user %d or book %d: %w", userID, bookID, core_errors.ErrNotFound)
		default:
			return core_domain.ShelfBook{}, fmt.Errorf("scan query row: %w", err)
		}
	}

	shelfDomain := core_domain.NewShelfBook(
		shelf.UserID,
		shelf.BookID,
		shelf.Version,
		shelf.Read,
		shelf.Rating,
		shelf.Review,
		shelf.AddedAt,
		shelf.ReadAt,
	)

	return shelfDomain, nil
}
