package statistics_repository_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/Doomed07/Bookshelf/internal/core/domain"
)

func (r *StatsRepository) GetBookShelves(ctx context.Context, bookID int) ([]core_domain.ShelfBook, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT user_id, book_id, version, read, rating, review, added_at, read_at, reviewed_at
	FROM bookshelfapp.bookshelf
	WHERE book_id = $1;
	`

	rows, err := r.pool.Query(ctx, query, bookID)
	if err != nil {
		return nil, fmt.Errorf("get query rows: %w", err)
	}
	defer rows.Close()

	var shelves []ShelfBookModel
	for rows.Next() {
		var s ShelfBookModel
		err := rows.Scan(
			&s.UserID,
			&s.BookID,
			&s.Version,
			&s.Read,
			&s.Rating,
			&s.Review,
			&s.AddedAt,
			&s.ReadAt,
			&s.ReviewedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan rows: %w", err)
		}
		shelves = append(shelves, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows next: %w", err)
	}

	shelvesDomain := shelvesDomainsFromModels(shelves)
	return shelvesDomain, nil
}
