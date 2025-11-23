package model

type PullRequestReviewer struct {
	ID         string `json:"pull_request_id"`
	ReviewerID string `json:"reviewer_id"`
}
