package users_repository_postgres

import (
	"time"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
)

type UserModel struct {
	ID        int
	Version   int
	Username  string
	Email     string
	CreatedAt time.Time
}

func userDomainsFromUserModels(users []UserModel) []domain.User {
	userDomains := make([]domain.User, len(users))

	for i, user := range users {
		userDomains[i] = domain.NewUser(
			user.ID,
			user.Version,
			user.Username,
			user.Email,
			user.CreatedAt,
		)
	}

	return userDomains
}
