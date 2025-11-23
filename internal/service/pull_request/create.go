package pull_request

import (
	"context"
	"math/rand"

	"pr-reviewer-assigment-service/internal/model"
)

func (s *serv) Create(ctx context.Context, pr *model.PullRequest) (*model.PullRequest, error) {
	createdPR, err := s.pullRequestRepository.Create(ctx, pr)
	if err != nil {
		return nil, err
	}

	author, err := s.userRepository.Get(ctx, pr.AuthorID)
	if err != nil {
		return nil, err
	}

	isActive := true
	activeMembers, err := s.userRepository.GetList(ctx, author.TeamName, &isActive)
	if err != nil {
		return nil, err
	}

	reviewers := selectRandomReviewers(activeMembers)

	for _, reviewer := range reviewers {
		if reviewer.ID == pr.AuthorID {
			continue
		}
		assignment := &model.PullRequestReviewer{
			ID:         createdPR.ID,
			ReviewerID: reviewer.ID,
		}
		_, err := s.pullRequestReviewersRepository.Create(ctx, assignment)
		if err != nil {
			return nil, err
		}
		createdPR.AssignedReviewers = append(createdPR.AssignedReviewers, reviewer.ID)
	}

	return createdPR, nil
}

func selectRandomReviewers(users []*model.User) []*model.User {
	if len(users) < 2 {
		return users
	}

	i := rand.Intn(len(users))
	j := rand.Intn(len(users))
	for j == i {
		j = rand.Intn(len(users))
	}

	return []*model.User{users[i], users[j]}
}
