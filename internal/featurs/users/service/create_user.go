package users_service

import (
	"context"
	"fmt"

	"github.com/Doomed07/Bookshelf/internal/core/core_domain"
)

func (s *UsersService) CreateUser(
	ctx context.Context,
	user core_domain.User,
) (core_domain.User, error) {
	if err := user.Validate(); err != nil {
		return core_domain.User{}, fmt.Errorf("validate user domain: %w", err)
	}

	user, err := s.usersRepository.CreateUser(ctx, user)
	if err != nil {
		return core_domain.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}
