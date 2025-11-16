package team

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrTeamNotFound = errors.New("NOT_FOUND")
)

func (s *Service) GetTeam(ctx context.Context, teamName string) (Team, error) {
	teamNameDB, repoUsers, err := s.repo.GetTeam(ctx, teamName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Team{}, ErrTeamNotFound
		}
		return Team{}, err
	}

	members := make([]TeamMember, len(repoUsers))
	for i, u := range repoUsers {
		members[i] = TeamMember{
			UserID:   u.UserID,
			Username: u.Username,
			IsActive: u.IsActive,
		}
	}

	return Team{
		TeamName: teamNameDB,
		Members:  members,
	}, nil
}
