package model

import (
	"time"
)

type PullRequest struct {
	ID                string    `json:"pull_request_id"`
	Name              string    `json:"pull_request_name"`
	AuthorID          string    `json:"author_id"`
	Status            string    `json:"status"`
	AssignedReviewers []string  `json:"assigned_reviewers"`
	MergedAt          time.Time `json:"mergedAt"`
}

type ReassignResult struct {
	PR            *PullRequest
	NewReviewerID string
}
