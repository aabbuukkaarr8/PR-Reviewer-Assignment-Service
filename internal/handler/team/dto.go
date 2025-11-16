package team

import (
	teamsrv "github.com/aabbuukkaarr8/PRService/internal/service/team"
)

type MemberTeam struct {
	UserID   string `json:"user_id" binding:"required"`
	Username string `json:"username" binding:"required"`
	IsActive bool   `json:"is_active"`
}

type Team struct {
	TeamName string       `json:"team_name" binding:"required"`
	Members  []MemberTeam `json:"members" binding:"required,dive"`
}

type CreateTeamRequest = Team

type CreateTeamResponse struct {
	Team Team `json:"team"`
}

type GetTeamResponse struct {
	Team Team `json:"team"`
}

func (t *Team) ToService() teamsrv.Team {
	members := make([]teamsrv.TeamMember, len(t.Members))
	for i, m := range t.Members {
		members[i] = teamsrv.TeamMember{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		}
	}

	return teamsrv.Team{
		TeamName: t.TeamName,
		Members:  members,
	}
}

func (t *Team) FillFromService(s teamsrv.Team) {
	members := make([]MemberTeam, len(s.Members))
	for i, m := range s.Members {
		members[i] = MemberTeam{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		}
	}

	t.TeamName = s.TeamName
	t.Members = members
}
