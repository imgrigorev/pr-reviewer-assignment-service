package pull_request

import (
	"context"
	"errors"

	"pr-reviewer-assigment-service/internal/model"
)

func (s *serv) Merge(ctx context.Context, prID string) (*model.PullRequest, error) {
	if prID == "" {
		return nil, errors.New("pull request ID is required")
	}

	mergedPR, err := s.pullRequestRepository.Merge(ctx, prID)
	if err != nil {
		return nil, err
	}

	reviewers, err := s.pullRequestReviewersRepository.GetByPrID(ctx, prID)
	if err != nil {
		return nil, err
	}

	for _, reviewer := range reviewers {
		mergedPR.AssignedReviewers = append(mergedPR.AssignedReviewers, reviewer.ReviewerID)
	}

	return mergedPR, nil
}
