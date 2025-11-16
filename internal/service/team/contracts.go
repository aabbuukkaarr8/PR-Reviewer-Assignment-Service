package team

import (
	"context"

	"github.com/aabbuukkaarr8/PRService/internal/repository/team"
)

type Repo interface {
	TeamExists(ctx context.Context, teamName string) (bool, error)
	CreateTeam(ctx context.Context, teamName string) error
	GetTeam(ctx context.Context, teamName string) (string, []team.User, error)
	UserExists(ctx context.Context, userID string) (bool, error)
	CreateUser(ctx context.Context, userID, username, teamName string, isActive bool) error
	UpdateUser(ctx context.Context, userID, username, teamName string, isActive bool) error
	BeginTx(ctx context.Context) (team.Tx, error)
}
