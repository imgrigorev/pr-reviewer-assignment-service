package converter

import (
	"pr-reviewer-assigment-service/internal/model"
	modelRepo "pr-reviewer-assigment-service/internal/repository/user/model"
)

func ToUserFromRepo(user *modelRepo.User) *model.User {
	return &model.User{
		ID:       user.ID,
		Name:     user.Name,
		IsActive: user.IsActive,
		TeamName: user.TeamName,
	}
}
