package pullrequest

import (
	"encoding/json"
	"errors"
	"net/http"

	"pr-reviewer-assigment-service/internal/model"
	service "pr-reviewer-assigment-service/internal/service/pull_request"
)

func (i *Implementation) ReassignReviewer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		PullRequestID string `json:"pull_request_id"`
		OldReviewerID string `json:"old_reviewer_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.PullRequestID == "" {
		http.Error(w, "pull_request_id is required", http.StatusBadRequest)
		return
	}
	if req.OldReviewerID == "" {
		http.Error(w, "old_reviewer_id is required", http.StatusBadRequest)
		return
	}

	result, err := i.pullRequestService.ReassignReviewer(ctx, req.PullRequestID, req.OldReviewerID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPRMerged):
			i.writeError(w, http.StatusConflict, "PR_MERGED", "cannot reassign on merged PR")
		case errors.Is(err, service.ErrReviewerNotAssigned):
			i.writeError(w, http.StatusConflict, "NOT_ASSIGNED", "reviewer is not assigned to this PR")
		case errors.Is(err, service.ErrReviewerNotFound):
			i.writeError(w, http.StatusNotFound, "REVIEWER_NOT_FOUND", "reviewer not found")
		case errors.Is(err, service.ErrNoCandidate):
			i.writeError(w, http.StatusConflict, "NO_CANDIDATE", "no active replacement candidate in team")
		case errors.Is(err, service.ErrPRNotFound):
			i.writeError(w, http.StatusNotFound, "PR_NOT_FOUND", "pull request not found")
		default:
			http.Error(w, "failed to reassign reviewer", http.StatusInternalServerError)
		}
		return
	}

	response := struct {
		PR         *model.PullRequest `json:"pr"`
		ReplacedBy string             `json:"replaced_by"`
	}{
		PR:         result.PR,
		ReplacedBy: result.NewReviewerID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (i *Implementation) writeError(w http.ResponseWriter, status int, code, message string) {
	errorResponse := struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}{}
	errorResponse.Error.Code = code
	errorResponse.Error.Message = message

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse)
}
