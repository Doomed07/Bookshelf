package bookshelf_repository_postgres

import (
	"context"
	"fmt"
)

func (r *BookshelfRepository) UserExists(ctx context.Context, userID int) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT EXISTS (SELECT 1 FROM bookshelfapp.users WHERE id = $1);
	`
	row := r.pool.QueryRow(ctx, query, userID)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, fmt.Errorf("scan query row: %w", err)
	}

	return exists, nil
}
