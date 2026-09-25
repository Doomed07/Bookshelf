package books_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
	core_errors "github.com/Doomed07/Bookshelf/internal/core/errors"
	core_postgres_pool "github.com/Doomed07/Bookshelf/internal/core/repository/postgres/pool"
)

func (r *BooksRepository) GetBook(ctx context.Context, id int) (domain.Book, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, title, author, year, pages, genres, description, score, reads_count
	FROM bookshelfapp.books
	WHERE id=$1;
	`

	row := r.pool.QueryRow(ctx, query, id)

	var b BookModel
	err := row.Scan(
		&b.ID,
		&b.Title,
		&b.Author,
		&b.Year,
		&b.Pages,
		&b.Genres,
		&b.Description,
		&b.Score,
		&b.ReadsCount,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Book{}, fmt.Errorf(
				"book with id %d: %w", id, core_errors.ErrNotFound)
		}
		return domain.Book{}, fmt.Errorf("failed to scan row: %w", err)
	}

	bookDomain := bookDomainFromModel(b)
	return bookDomain, nil
}
