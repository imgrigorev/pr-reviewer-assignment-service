package model

import "time"

type PullRequestReviewer struct {
	PrID       string    `json:"pr_id"`
	ReviewerID string    `json:"reviewer_id"`
	AssignedAt time.Time `json:"assigned_at"`
}
