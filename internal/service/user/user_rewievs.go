package user

import (
	"context"
	"errors"

	"pr-reviewer-assigment-service/internal/model"

	"github.com/jackc/pgx/v5"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

func (s *serv) GetUserReviews(ctx context.Context, userID string) (*model.UserReviewsResponse, error) {
	user, err := s.userRepository.Get(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("user is not active")
	}

	prs, err := s.pullRequestRepository.GetByReviewerID(ctx, userID)
	if err != nil {
		return nil, err
	}

	shortPRs := make([]model.PullRequestShort, 0, len(prs))
	for _, pr := range prs {
		shortPRs = append(shortPRs, model.PullRequestShort{
			ID:       pr.ID,
			Name:     pr.Name,
			AuthorID: pr.AuthorID,
			Status:   pr.Status,
		})
	}

	return &model.UserReviewsResponse{
		UserID:       userID,
		PullRequests: shortPRs,
	}, nil
}
