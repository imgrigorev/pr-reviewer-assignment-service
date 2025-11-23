package converter

import (
	"pr-reviewer-assigment-service/internal/model"
	modelRepo "pr-reviewer-assigment-service/internal/repository/pull_request_reviewers/model"
)

func ToPullRequestReviewerFromRepo(repo *modelRepo.PullRequestReviewer) *model.PullRequestReviewer {
	if repo == nil {
		return nil
	}

	return &model.PullRequestReviewer{
		ID:         repo.PrID,
		ReviewerID: repo.ReviewerID,
	}
}

func ToPullRequestReviewersFromRepo(repo []*modelRepo.PullRequestReviewer) []*model.PullRequestReviewer {
	if repo == nil {
		return nil
	}

	result := make([]*model.PullRequestReviewer, 0, len(repo))
	for _, r := range repo {
		result = append(result, ToPullRequestReviewerFromRepo(r))
	}

	return result
}
