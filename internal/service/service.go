package service

import (
	"context"

	"pr-reviewer-assigment-service/internal/model"
)

type UserService interface {
	Create(ctx context.Context, u *model.User, teamName string) (*model.User, error)
	Get(ctx context.Context, id string) (*model.User, error)
	Update(ctx context.Context, u *model.User) (*model.User, error)
	GetUserReviews(ctx context.Context, userID string) (*model.UserReviewsResponse, error)
}

type TeamService interface {
	Create(ctx context.Context, team *model.Team) (*model.Team, error)
	Get(ctx context.Context, teamName string, isActive *bool) (*model.Team, error)
}

type PullRequestService interface {
	Create(ctx context.Context, pr *model.PullRequest) (*model.PullRequest, error)
	Merge(ctx context.Context, prID string) (*model.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID, oldReviewerID string) (*model.ReassignResult, error)
}
