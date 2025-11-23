package pullrequest

import "pr-reviewer-assigment-service/internal/service"

type Implementation struct {
	pullRequestService service.PullRequestService
}

func NewImplementation(pullRequestService service.PullRequestService) *Implementation {
	return &Implementation{
		pullRequestService: pullRequestService,
	}
}
