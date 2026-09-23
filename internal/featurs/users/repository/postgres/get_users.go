package users_repository_postgres

import (
	"context"
	"fmt"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
)

func (r *UsersRepository) GetUsers(
	ctx context.Context,
	limit, offset *int,
) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, username, email, created_at
	FROM bookshelfapp.users
	ORDER BY id ASC
	LIMIT $1
	OFFSET $2
 	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get query rows: %w", err)
	}
	defer rows.Close()

	var userModels []UserModel
	for rows.Next() {
		var userModel UserModel
		err := rows.Scan(
			&userModel.ID,
			&userModel.Version,
			&userModel.Username,
			&userModel.Email,
			&userModel.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan query row: %w", err)
		}

		userModels = append(userModels, userModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	userDomains := userDomainsFromUserModels(userModels)

	return userDomains, nil
}
