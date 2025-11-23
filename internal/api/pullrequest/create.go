package pullrequest

import (
	"encoding/json"
	"errors"
	"net/http"

	innerError "pr-reviewer-assigment-service/internal/errors"
	"pr-reviewer-assigment-service/internal/model"
)

type CreatePullRequestRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type CreatePullRequestResponse struct {
	PR *PullRequestResponse `json:"pr"`
}

type PullRequestResponse struct {
	ID                string   `json:"pull_request_id"`
	Name              string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
}

func (i *Implementation) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreatePullRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.PullRequestID == "" || req.PullRequestName == "" || req.AuthorID == "" {
		http.Error(w, "all fields are required", http.StatusBadRequest)
		return
	}

	pr := &model.PullRequest{
		ID:       req.PullRequestID,
		Name:     req.PullRequestName,
		AuthorID: req.AuthorID,
		Status:   "OPEN",
	}

	createdPR, err := i.pullRequestService.Create(ctx, pr)
	if err != nil {
		if errors.Is(err, innerError.ErrPRAlreadyExists) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]any{
				"error": map[string]string{
					"code":    "PR_EXISTS",
					"message": "PR id already exists",
				},
			})
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := CreatePullRequestResponse{
		PR: &PullRequestResponse{
			ID:                createdPR.ID,
			Name:              createdPR.Name,
			AuthorID:          createdPR.AuthorID,
			Status:            createdPR.Status,
			AssignedReviewers: createdPR.AssignedReviewers,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
