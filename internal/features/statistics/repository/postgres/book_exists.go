package statistics_repository_postgres

import (
	"context"
	"fmt"
)

func (r *StatsRepository) BookExists(ctx context.Context, bookID int) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT EXISTS (SELECT 1 FROM bookshelfapp.books WHERE id = $1);
	`

	row := r.pool.QueryRow(ctx, query, bookID)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, fmt.Errorf("row scan: %w", err)
	}

	return exists, nil
}
