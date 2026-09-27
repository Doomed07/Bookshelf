package bookshelf_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
	core_errors "github.com/Doomed07/Bookshelf/internal/core/errors"
	core_postgres_pool "github.com/Doomed07/Bookshelf/internal/core/repository/postgres/pool"
)

func (r *BookshelfRepository) GetBook(ctx context.Context, userID, bookID int) (domain.ShelfBookWithBook, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT s.user_id, s.book_id, s.version, s.read, s.rating, s.review, s.added_at, s.read_at, 
	b.id, b.title, b.author, b.year, b.pages, b.genres, b.description, b.score, b.reads_count
	FROM bookshelfapp.bookshelf s
	JOIN bookshelfapp.books b ON b.id = s.book_id
	WHERE s.user_id = $1 AND s.book_id = $2;
	`
	row := r.pool.QueryRow(ctx, query, userID, bookID)

	var shelf ShelfBookWithBookModel
	err := row.Scan(
		&shelf.ShelfBook.UserID,
		&shelf.ShelfBook.BookID,
		&shelf.ShelfBook.Version,
		&shelf.ShelfBook.Read,
		&shelf.ShelfBook.Rating,
		&shelf.ShelfBook.Review,
		&shelf.ShelfBook.AddedAt,
		&shelf.ShelfBook.ReadAt,
		&shelf.Book.ID,
		&shelf.Book.Title,
		&shelf.Book.Author,
		&shelf.Book.Year,
		&shelf.Book.Pages,
		&shelf.Book.Genres,
		&shelf.Book.Description,
		&shelf.Book.Score,
		&shelf.Book.ReadsCount,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.ShelfBookWithBook{}, fmt.Errorf(
				"book %d on bookshelf of user %d: %w", bookID, userID, core_errors.ErrNotFound)
		}
		return domain.ShelfBookWithBook{}, fmt.Errorf("scan query row: %w", err)
	}

	SBWBDomain := domainSBWBFromModel(shelf)
	return SBWBDomain, nil
}
