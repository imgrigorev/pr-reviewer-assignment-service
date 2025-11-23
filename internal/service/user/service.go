package user

import (
	"pr-reviewer-assigment-service/internal/repository"
	"pr-reviewer-assigment-service/internal/service"
)

type serv struct {
	userRepository        repository.UserRepository
	pullRequestRepository repository.PullRequestRepository
}

func NewService(userRepository repository.UserRepository,
	pullRequestRepository repository.PullRequestRepository) service.UserService {
	return &serv{
		userRepository:        userRepository,
		pullRequestRepository: pullRequestRepository,
	}
}
