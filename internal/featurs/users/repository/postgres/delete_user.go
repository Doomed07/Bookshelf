package users_repository_postgres

import (
	"context"
	"fmt"

	core_errors "github.com/Doomed07/Bookshelf/internal/core/errors"
)

func (r *UsersRepository) DeleteUser(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM bookshelfapp.users
	WHERE id=$1
	`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user with id: %d: %w",
			id, core_errors.ErrNotFound)
	}

	return nil
}
