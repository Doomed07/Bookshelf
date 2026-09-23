package users_service

import (
	"context"
	"fmt"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
	core_errors "github.com/Doomed07/Bookshelf/internal/core/errors"
)

func (s *UsersService) GetUsers(
	ctx context.Context,
	limit, offset *int,
) ([]domain.User, error) {
	if limit != nil && (*limit < 1 || *limit > 100) {
		return nil, fmt.Errorf("invalid 'limit' query param: %w", core_errors.ErrInvalidArgument)
	}

	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf("invalid 'offset' query param: %w", core_errors.ErrInvalidArgument)
	}

	users, err := s.usersRepository.GetUsers(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get users from repo: %w", err)
	}

	return users, nil
}
