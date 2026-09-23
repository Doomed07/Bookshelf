package users_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
	core_errors "github.com/Doomed07/Bookshelf/internal/core/errors"
	"github.com/jackc/pgx/v5"
)

func (r *UsersRepository) PatchUser(ctx context.Context, id int, user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE bookshelfapp.users
	SET username=$1, email=$2, version=version+1
	WHERE id=$3 AND version=$4
	RETURNING id, version, username, email, created_at
	`

	row := r.pool.QueryRow(ctx, query, user.Username, user.Email, id, user.Version)

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.Username,
		&userModel.Email,
		&userModel.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"user with id: %d - cocurrently accesed: %w",
				id, core_errors.ErrConflict,
			)
		}
		return domain.User{}, fmt.Errorf("scan query row: %w", err)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.Username,
		userModel.Email,
		userModel.CreatedAt,
	)

	return userDomain, nil
}
