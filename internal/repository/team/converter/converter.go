package converter

import (
	"pr-reviewer-assigment-service/internal/model"
	modelRepo "pr-reviewer-assigment-service/internal/repository/team/model"
)

func ToTeamFromRepo(user *modelRepo.Team) *model.Team {
	return &model.Team{
		Name: user.Name,
	}
}
