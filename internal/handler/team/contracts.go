package team

import (
	"context"

	teamsrv "github.com/aabbuukkaarr8/PRService/internal/service/team"
)

type ServiceTeam interface {
	CreateTeam(ctx context.Context, team teamsrv.Team) (teamsrv.Team, error)
	GetTeam(ctx context.Context, teamName string) (teamsrv.Team, error)
}
