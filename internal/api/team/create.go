package team

import (
	"encoding/json"
	"errors"
	"net/http"

	innerError "pr-reviewer-assigment-service/internal/errors"
	"pr-reviewer-assigment-service/internal/model"
)

type CreateTeamResponse struct {
	Team *TeamResponse `json:"team"`
}

type TeamResponse struct {
	TeamName string           `json:"team_name"`
	Members  []MemberResponse `json:"members"`
}

type MemberResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

func (i *Implementation) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.Team
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	team, err := i.teamService.Create(ctx, &req)
	if err != nil {
		if errors.Is(err, innerError.ErrTeamExists) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]any{
				"error": map[string]string{
					"code":    "TEAM_EXISTS",
					"message": "team_name already exists",
				},
			})
			return
		}

		http.Error(w, "failed to create team", http.StatusInternalServerError)
		return
	}

	response := i.toTeamResponse(team)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (i *Implementation) toTeamResponse(team *model.Team) CreateTeamResponse {
	members := make([]MemberResponse, 0, len(team.Members))
	for _, member := range team.Members {
		members = append(members, MemberResponse{
			UserID:   member.ID,
			Username: member.Name,
			IsActive: member.IsActive,
		})
	}

	return CreateTeamResponse{
		Team: &TeamResponse{
			TeamName: team.Name,
			Members:  members,
		},
	}
}
