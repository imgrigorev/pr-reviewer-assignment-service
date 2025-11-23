package repository

import (
	"context"

	"pr-reviewer-assigment-service/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, u *model.User, teamName string) (*model.User, error)
	Get(ctx context.Context, id string) (*model.User, error)
	Update(ctx context.Context, u *model.User) (*model.User, error)
	GetList(ctx context.Context, teamName string, isActive *bool) ([]*model.User, error)
}

type TeamRepository interface {
	Create(ctx context.Context, name string) (*model.Team, error)
	GetByName(ctx context.Context, name string) (*model.Team, error)
}

type PullRequestRepository interface {
	Create(ctx context.Context, pr *model.PullRequest) (*model.PullRequest, error)
	Merge(ctx context.Context, prID string) (*model.PullRequest, error)
	GetByID(ctx context.Context, id string) (*model.PullRequest, error)
	GetByReviewerID(ctx context.Context, reviewerID string) ([]*model.PullRequest, error)
}

type PullRequestReviewersRepository interface {
	Create(ctx context.Context, assignment *model.PullRequestReviewer) (*model.PullRequestReviewer, error)
	GetByUserID(ctx context.Context, userID string) ([]*model.PullRequestReviewer, error)
	GetByPrID(ctx context.Context, prID string) ([]*model.PullRequestReviewer, error)
	Delete(ctx context.Context, prID, reviewerID string) error
	DeleteByPrID(ctx context.Context, prID string) error
	DeleteByUserID(ctx context.Context, userID string) error
}
