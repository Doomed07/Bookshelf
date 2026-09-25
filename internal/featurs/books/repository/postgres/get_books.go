package books_repository_postgres

import (
	"context"
	"fmt"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
)

func (r *BooksRepository) GetBooks(
	ctx context.Context,
	title, author *string,
	limit, offset int,
) ([]domain.Book, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, title, author, year, pages, genres, description, score, reads_count
	FROM bookshelfapp.books
	WHERE ($1::text IS NULL OR title ILIKE $1)
	  AND ($2::text IS NULL OR author ILIKE $2)
	ORDER BY id ASC
	LIMIT $3
	OFFSET $4;
	`

	rows, err := r.pool.Query(
		ctx,
		query,
		likePattern(title),
		likePattern(author),
		limit,
		offset)
	if err != nil {
		return nil, fmt.Errorf("get query rows: %w", err)
	}
	defer rows.Close()

	var booksModel []BookModel
	for rows.Next() {
		var b BookModel
		err := rows.Scan(
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
			return nil, fmt.Errorf("scan row: %w", err)
		}
		booksModel = append(booksModel, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next row: %w", err)
	}

	domainBooks := booksDomainsFromModels(booksModel)

	return domainBooks, nil
}
