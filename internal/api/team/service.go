package team

import "pr-reviewer-assigment-service/internal/service"

type Implementation struct {
	teamService service.TeamService
}

func NewImplementation(teamService service.TeamService) *Implementation {
	return &Implementation{
		teamService: teamService,
	}
}
