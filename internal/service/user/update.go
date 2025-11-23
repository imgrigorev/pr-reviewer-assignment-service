package user

import (
	"context"

	"pr-reviewer-assigment-service/internal/model"
)

func (s *serv) Update(ctx context.Context, u *model.User) (*model.User, error) {
	user, err := s.userRepository.Update(ctx, u)
	if err != nil {
		return nil, err
	}

	return user, nil
}
