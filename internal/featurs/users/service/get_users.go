package users_service

import (
	"context"
	"fmt"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
)

func (s *UsersService) GetUsers(
	ctx context.Context,
	limit, offset *int,
) ([]domain.User, error) {
	lim, off, err := domain.NormalizePagination(limit, offset)
	if err != nil {
		return nil, err
	}

	users, err := s.usersRepository.GetUsers(ctx, lim, off)
	if err != nil {
		return nil, fmt.Errorf("get users from repo: %w", err)
	}

	return users, nil
}
