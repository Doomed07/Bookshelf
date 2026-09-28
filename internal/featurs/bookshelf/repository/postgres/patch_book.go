package bookshelf_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Doomed07/Bookshelf/internal/core/core_domain"
	core_errors "github.com/Doomed07/Bookshelf/internal/core/errors"
	core_postgres_pool "github.com/Doomed07/Bookshelf/internal/core/repository/postgres/pool"
)

func (r *BookshelfRepository) PatchBook(
	ctx context.Context,
	shelfBook core_domain.ShelfBook,
) (core_domain.ShelfBook, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE bookshelfapp.bookshelf s
	SET read = $1, rating = $2, review = $3, read_at = $4, version = version + 1
	WHERE user_id = $5 AND book_id = $6 AND version = $7
	RETURNING user_id, book_id, version, read, rating, review, added_at, read_at;
	`
	row := r.pool.QueryRow(
		ctx,
		query,
		shelfBook.Read,
		shelfBook.Rating,
		shelfBook.Review,
		shelfBook.ReadAt,
		shelfBook.UserID,
		shelfBook.BookID,
		shelfBook.Version)

	var s ShelfBookModel
	err := row.Scan(
		&s.UserID,
		&s.BookID,
		&s.Version,
		&s.Read,
		&s.Rating,
		&s.Review,
		&s.AddedAt,
		&s.ReadAt,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return core_domain.ShelfBook{}, fmt.Errorf(
				"book %d on bookshelf of user %d was modified concurrently: %w",
				shelfBook.BookID, shelfBook.UserID, core_errors.ErrConflict)
		}
		return core_domain.ShelfBook{}, fmt.Errorf("scan query row: %w", err)

	}

	shelfDomain := core_domain.NewShelfBook(
		s.UserID,
		s.BookID,
		s.Version,
		s.Read,
		s.Rating,
		s.Review,
		s.AddedAt,
		s.ReadAt,
	)

	return shelfDomain, nil
}
