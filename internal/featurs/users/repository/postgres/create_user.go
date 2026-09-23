package users_repository_postgres

import (
	"context"
	"fmt"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
)

func (r *UsersRepository) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO bookshelfapp.users (username, email)
	VALUES ($1, $2)
	RETURNING id, version, username, email, created_at
	`
	row := r.pool.QueryRow(ctx, query, user.Username, user.Email)

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.Username,
		&userModel.Email,
		&userModel.CreatedAt,
	)
	if err != nil {
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
