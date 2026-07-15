package users_transport_http

import "github.com/catstyle1101/todo_app_go/cmd/internal/core/domain"

type UserDTOResponse struct {
	ID          int     `json:"id" example:"999"`
	Version     int     `json:"version" example:"2"`
	FullName    string  `json:"full_name" example:"Ivan Ivanov"`
	PhoneNumber *string `json:"phone_number" example:"+79998887766"`
}

func userDTOFromDomain(domain domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:          domain.ID,
		Version:     domain.Version,
		FullName:    domain.FullName,
		PhoneNumber: domain.PhoneNumber,
	}
}

func usersDTOFromDomain(users []domain.User) []UserDTOResponse {
	usersDTO := make([]UserDTOResponse, len(users))
	for i, user := range users {
		usersDTO[i] = userDTOFromDomain(user)
	}

	return usersDTO
}
