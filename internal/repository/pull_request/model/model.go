package model

import (
	"database/sql"
	"time"
)

type PullRequest struct {
	ID                string       `json:"pull_request_id"`
	Name              string       `json:"pull_request_name"`
	AuthorID          string       `json:"author_id"`
	Status            string       `json:"status"`
	AssignedReviewers []string     `json:"assigned_reviewers"`
	CreatedAt         time.Time    `json:"created_at"`
	MergedAt          sql.NullTime `json:"merged_at"`
}
