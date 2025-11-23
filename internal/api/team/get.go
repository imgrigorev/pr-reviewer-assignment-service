package team

import (
	"encoding/json"
	"errors"
	"net/http"

	innerError "pr-reviewer-assigment-service/internal/errors"
)

//type GetTeamResponse struct {
//	Team *TeamResponse `json:"team"`
//}

func (i *Implementation) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	teamName := r.URL.Query().Get("team_name")

	if teamName == "" {
		http.Error(w, "team_name is required", http.StatusBadRequest)
		return
	}

	team, err := i.teamService.Get(ctx, teamName, nil)
	if err != nil {
		if errors.Is(err, innerError.ErrTeamNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]any{
				"error": map[string]string{
					"code":    "TEAM_NOT_FOUND",
					"message": "team not found",
				},
			})
			return
		}

		http.Error(w, "failed to get team", http.StatusInternalServerError)
		return
	}

	response := i.toTeamResponse(team)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

//func (i *Implementation) toTeamResponseDTO(team *model.Team) *TeamResponse {
//	if team == nil {
//		return nil
//	}
//
//	members := make([]MemberResponse, 0, len(team.Members))
//	for _, member := range team.Members {
//		members = append(members, MemberResponse{
//			UserID:   member.ID,
//			Username: member.Name,
//			IsActive: member.IsActive,
//		})
//	}
//
//	return &TeamResponse{
//		TeamName: team.Name,
//		Members:  members,
//	}
//}
