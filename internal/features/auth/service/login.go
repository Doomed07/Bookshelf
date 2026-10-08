package auth_service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	core_domain "github.com/Doomed07/Bookshelf/internal/core/domain"
	core_errors "github.com/Doomed07/Bookshelf/internal/core/errors"
	"golang.org/x/crypto/bcrypt"
)

var errInvalidCredentials = fmt.Errorf("invalid login or password: %w", core_errors.ErrUnauthorized)

func (s *AuthService) Login(ctx context.Context, login, password string) (core_domain.User, string, error) {
	login = strings.TrimSpace(login)
	if len(login) > 254 || len(password) > core_domain.MaxPasswordLen {
		return core_domain.User{}, "", errInvalidCredentials
	}

	creds, err := s.repo.GetCredentialsByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			_ = s.comparePassword(ctx, s.dummyHash, password)
			return core_domain.User{}, "", errInvalidCredentials
		}
		return core_domain.User{}, "", fmt.Errorf("get credentials: %w", err)
	}

	if err := s.comparePassword(ctx, []byte(creds.PasswordHash), password); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return core_domain.User{}, "", errInvalidCredentials
		}
		return core_domain.User{}, "", fmt.Errorf("compare password: %w", err)
	}

	_ = s.repo.DeleteExpiredSessions(ctx, s.now())
	token, err := s.startSession(ctx, creds.User.ID)
	if err != nil {
		return core_domain.User{}, "", err
	}
	return creds.User, token, nil
}
