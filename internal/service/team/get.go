package team

import (
	"context"

	"pr-reviewer-assigment-service/internal/model"
)

func (s *serv) Get(ctx context.Context, teamName string, isActive *bool) (*model.Team, error) {
	m, err := s.teamRepository.GetByName(ctx, teamName)
	if err != nil {
		return nil, err
	}

	users, err := s.userRepository.GetList(ctx, teamName, isActive)
	if err != nil {
		return nil, err
	}

	m.Members = users
	return m, nil
}
