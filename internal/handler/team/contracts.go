package team

import (
	"context"

	"github.com/aabbuukkaarr8/PRService/internal/service/team"
)

// TeamService интерфейс для работы с командами
type TeamService interface {
	CreateTeam(ctx context.Context, team team.Team) (team.Team, error)
	GetTeam(ctx context.Context, teamName string) (team.Team, error)
}
