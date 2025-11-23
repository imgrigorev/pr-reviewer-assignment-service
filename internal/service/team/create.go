package team

import (
	"context"

	"pr-reviewer-assigment-service/internal/model"
)

func (s *serv) Create(ctx context.Context, team *model.Team) (*model.Team, error) {

	t, err := s.teamRepository.Create(ctx, team.Name)
	if err != nil {
		return nil, err
	}

	for _, user := range team.Members {
		member, err := s.userRepository.Create(ctx, user, t.Name)
		if err != nil {
			return nil, err
		}
		t.Members = append(t.Members, member)
	}

	return t, nil
}
