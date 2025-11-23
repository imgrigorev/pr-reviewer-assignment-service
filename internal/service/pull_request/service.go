package pull_request

import (
	"pr-reviewer-assigment-service/internal/repository"
	"pr-reviewer-assigment-service/internal/service"
)

type serv struct {
	pullRequestRepository          repository.PullRequestRepository
	pullRequestReviewersRepository repository.PullRequestReviewersRepository
	teamRepository                 repository.TeamRepository
	userRepository                 repository.UserRepository
}

func NewService(pullRequestRepository repository.PullRequestRepository,
	pullRequestReviewersRepository repository.PullRequestReviewersRepository,
	teamRepository repository.TeamRepository,
	userRepository repository.UserRepository,
) service.PullRequestService {
	return &serv{
		pullRequestRepository:          pullRequestRepository,
		pullRequestReviewersRepository: pullRequestReviewersRepository,
		teamRepository:                 teamRepository,
		userRepository:                 userRepository,
	}
}
