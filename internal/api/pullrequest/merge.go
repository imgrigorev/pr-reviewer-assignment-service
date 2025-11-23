package pullrequest

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"pr-reviewer-assigment-service/internal/model"
)

func (i *Implementation) Merge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		PullRequestID string `json:"pull_request_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.PullRequestID == "" {
		http.Error(w, "pull_request_id is required", http.StatusBadRequest)
		return
	}

	mergedPR, err := i.pullRequestService.Merge(ctx, req.PullRequestID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "not found"):
			http.Error(w, "pull request not found", http.StatusNotFound)
		default:
			http.Error(w, "failed to merge pull request", http.StatusInternalServerError)
		}
		return
	}

	response := struct {
		PR *model.PullRequest `json:"pr"`
	}{
		PR: mergedPR,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
