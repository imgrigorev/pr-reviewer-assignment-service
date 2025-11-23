package user

import (
	"encoding/json"
	"errors"
	"net/http"

	innerError "pr-reviewer-assigment-service/internal/errors"
	"pr-reviewer-assigment-service/internal/model"
)

type createUserRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

func (i *Implementation) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	m := model.User{ID: req.UserID, IsActive: req.IsActive}

	createdUser, err := i.userService.Update(ctx, &m)
	if err != nil {
		if errors.Is(err, innerError.ErrUserNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]any{
				"error": map[string]string{
					"code":    "USER_NOT_FOUND",
					"message": "user not found",
				},
			})
			return
		}
		http.Error(w, "failed to update user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(createdUser)
}
