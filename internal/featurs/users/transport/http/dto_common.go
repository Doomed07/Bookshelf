package users_transport_http

import "github.com/Doomed07/Bookshelf/internal/core/domain"

type UserDTOResponse struct {
	ID       int    `json:"id"`
	Version  int    `json:"version"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func userDTOFromDomain(user domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:       user.ID,
		Version:  user.Version,
		Username: user.Username,
		Email:    user.Email,
	}
}

func usersDTOFromDomains(users []domain.User) []UserDTOResponse {
	usersDTO := make([]UserDTOResponse, len(users))

	for i, v := range users {
		usersDTO[i] = userDTOFromDomain(v)
	}

	return usersDTO
}