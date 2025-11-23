package pull_request

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"time"

	"pr-reviewer-assigment-service/internal/model"
)

var (
	ErrPRMerged            = errors.New("cannot reassign on merged PR")
	ErrPRNotFound          = errors.New("pull request not found")
	ErrNoCandidate         = errors.New("no active replacement candidate in team")
	ErrReviewerNotAssigned = errors.New("reviewer not assigned")
	ErrReviewerNotFound    = errors.New("reviewer not found")
)

func (s *serv) ReassignReviewer(ctx context.Context, prID, oldReviewerID string) (*model.ReassignResult, error) {
	pr, err := s.pullRequestRepository.GetByID(ctx, prID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPRNotFound
		}
		return nil, err
	}

	if pr.Status == "MERGED" {
		return nil, ErrPRMerged
	}

	reviewers, err := s.pullRequestReviewersRepository.GetByPrID(ctx, prID)
	if err != nil {
		return nil, err
	}

	found := false
	for _, reviewer := range reviewers {
		if reviewer.ReviewerID == oldReviewerID {
			found = true
			break
		}
	}
	if !found {
		return nil, ErrReviewerNotAssigned
	}

	oldReviewer, err := s.userRepository.Get(ctx, oldReviewerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrReviewerNotFound
		}
		return nil, err
	}

	isActive := true
	activeMembers, err := s.userRepository.GetList(ctx, oldReviewer.TeamName, &isActive)
	if err != nil {
		return nil, err
	}

	availableReviewers := excludeCurrentReviewers(activeMembers, reviewers, oldReviewerID, pr.AuthorID)

	if len(availableReviewers) == 0 {
		return nil, ErrNoCandidate
	}

	newReviewer := selectRandomReviewer(availableReviewers)

	err = s.pullRequestReviewersRepository.Delete(ctx, prID, oldReviewerID)
	if err != nil {
		return nil, err
	}

	newAssignment := &model.PullRequestReviewer{
		ID:         prID,
		ReviewerID: newReviewer.ID,
	}
	_, err = s.pullRequestReviewersRepository.Create(ctx, newAssignment)
	if err != nil {
		return nil, err
	}

	updatedPR, err := s.pullRequestRepository.GetByID(ctx, prID)
	if err != nil {
		return nil, err
	}

	updatedReviewers, err := s.pullRequestReviewersRepository.GetByPrID(ctx, prID)
	if err != nil {
		return nil, err
	}

	for _, reviewer := range updatedReviewers {
		updatedPR.AssignedReviewers = append(updatedPR.AssignedReviewers, reviewer.ReviewerID)
	}

	return &model.ReassignResult{
		PR:            updatedPR,
		NewReviewerID: newReviewer.ID,
	}, nil
}

func excludeCurrentReviewers(allUsers []*model.User, currentReviewers []*model.PullRequestReviewer, excludeUserID string, authorID string) []*model.User {
	currentReviewerIDs := make(map[string]bool)
	for _, reviewer := range currentReviewers {
		currentReviewerIDs[reviewer.ReviewerID] = true
	}

	var available []*model.User
	for _, user := range allUsers {
		if !currentReviewerIDs[user.ID] && user.ID != excludeUserID && user.ID != authorID {
			available = append(available, user)
		}
	}

	return available
}

func selectRandomReviewer(users []*model.User) *model.User {
	if len(users) == 0 {
		return nil
	}
	rand.Seed(time.Now().UnixNano())
	return users[rand.Intn(len(users))]
}
