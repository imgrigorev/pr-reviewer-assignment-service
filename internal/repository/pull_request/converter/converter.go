package converter

import (
	"time"

	"pr-reviewer-assigment-service/internal/model"
	modelRepo "pr-reviewer-assigment-service/internal/repository/pull_request/model"
)

func ToPullRequestFromRepo(pr *modelRepo.PullRequest) *model.PullRequest {
	var mergedAt time.Time
	if pr.MergedAt.Valid {
		mergedAt = pr.MergedAt.Time
	}
	return &model.PullRequest{
		ID:                pr.ID,
		Name:              pr.Name,
		AuthorID:          pr.AuthorID,
		Status:            pr.Status,
		AssignedReviewers: nil,
		MergedAt:          mergedAt,
	}
}
