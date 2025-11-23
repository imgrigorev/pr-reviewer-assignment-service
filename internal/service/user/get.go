package user

import (
	"context"

	"pr-reviewer-assigment-service/internal/model"
)

func (s *serv) Get(ctx context.Context, id string) (*model.User, error) {
	user, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
