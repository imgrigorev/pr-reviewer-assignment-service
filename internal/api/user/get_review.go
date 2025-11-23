package user

import (
	"encoding/json"
	"errors"
	"net/http"

	service "pr-reviewer-assigment-service/internal/service/user"
)

func (i *Implementation) GetUserReviews(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		i.writeError(w, http.StatusBadRequest, "MISSING_PARAMETER", "user_id query parameter is required")
		return
	}

	response, err := i.userService.GetUserReviews(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			i.writeError(w, http.StatusNotFound, "USER_NOT_FOUND", "user not found")
		default:
			i.writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to get user reviews")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		i.writeError(w, http.StatusInternalServerError, "ENCODING_ERROR", "failed to encode response")
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
