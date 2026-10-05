package statistics_repository_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/Doomed07/Bookshelf/internal/core/domain"
)

func (r *StatsRepository) GetTopBooksByScore(ctx context.Context, minRatings, top int) ([]core_domain.Book, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, title, author, year, pages, genres, description, score, reads_count
	FROM bookshelfapp.books
	WHERE rating_count >= $1
	ORDER BY score DESC NULLS LAST, id
	LIMIT $2;
	`

	rows, err := r.pool.Query(ctx, query, minRatings, top)
	if err != nil {
		return nil, fmt.Errorf("get query rows: %w", err)
	}
	defer rows.Close()

	var books []BookModel
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
		books = append(books, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	booksDomain := booksDomainsFromModels(books)

	return booksDomain, nil
}
