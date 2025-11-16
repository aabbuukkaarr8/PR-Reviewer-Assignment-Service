package team

import (
	"context"
	"errors"
)

var (
	ErrTeamExists = errors.New("TEAM_EXISTS")
)

func (s *Service) CreateTeam(ctx context.Context, team Team) (Team, error) {
	exists, err := s.repo.TeamExists(ctx, team.TeamName)
	if err != nil {
		return Team{}, err
	}
	if exists {
		return Team{}, ErrTeamExists
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return Team{}, err
	}
	defer tx.Rollback()

	if err := tx.CreateTeam(team.TeamName); err != nil {
		return Team{}, err
	}

	for _, member := range team.Members {
		userExists, err := s.repo.UserExists(ctx, member.UserID)
		if err != nil {
			return Team{}, err
		}

		if userExists {
			if err := tx.UpdateUser(member.UserID, member.Username, team.TeamName, member.IsActive); err != nil {
				return Team{}, err
			}
		} else {
			if err := tx.CreateUser(member.UserID, member.Username, team.TeamName, member.IsActive); err != nil {
				return Team{}, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return Team{}, err
	}

	return team, nil
}
