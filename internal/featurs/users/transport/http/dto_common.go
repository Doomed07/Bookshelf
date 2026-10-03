package users_transport_http

import core_domain "github.com/Doomed07/Bookshelf/internal/core/domain"

type UserDTOResponse struct {
	ID       int    `json:"id" example:"1"`
	Version  int    `json:"version" example:"3"`
	Username string `json:"username" example:"book_worm07"`
	Email    string `json:"email" example:"lol@mail.com"`
}

func userDTOFromDomain(user core_domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:       user.ID,
		Version:  user.Version,
		Username: user.Username,
		Email:    user.Email,
	}
}

func usersDTOFromDomains(users []core_domain.User) []UserDTOResponse {
	usersDTO := make([]UserDTOResponse, len(users))

	for i, v := range users {
		usersDTO[i] = userDTOFromDomain(v)
	}

	return usersDTO
}
