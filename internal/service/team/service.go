package team

import (
	"pr-reviewer-assigment-service/internal/repository"
	"pr-reviewer-assigment-service/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
	teamRepository repository.TeamRepository
}

func NewService(userRepository repository.UserRepository, teamRepository repository.TeamRepository) service.TeamService {
	return &serv{
		userRepository: userRepository,
		teamRepository: teamRepository,
	}
}
