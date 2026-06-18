package users_service

import (
	"context"
	"fmt"

	"github.com/catstyle1101/todo_app_go/cmd/internal/core/domain"
)

func (s *UsersService) GetUser(
	ctx context.Context,
	userID int,
) (domain.User, error) {
	user, err := s.usersRepository.GetUser(ctx, userID)

	if err != nil {
		return domain.User{}, fmt.Errorf("get user from repository: %w", err)
	}

	return user, nil
}
