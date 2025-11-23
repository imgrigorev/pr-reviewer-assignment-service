package user

import (
	"context"

	"pr-reviewer-assigment-service/internal/model"
)

func (s *serv) Create(ctx context.Context, u *model.User, teamName string) (*model.User, error) {
	user, err := s.userRepository.Create(ctx, u, teamName)
	if err != nil {
		return nil, err
	}

	return user, nil
}
