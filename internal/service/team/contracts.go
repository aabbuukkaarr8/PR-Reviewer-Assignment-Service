package team

import (
	"context"

	"github.com/aabbuukkaarr8/PRService/internal/repository/team"
)

// RepoUser структура пользователя из repository (для передачи данных между repository и service)
type RepoUser struct {
	UserID   string
	Username string
	IsActive bool
}

type Repo interface {
	// TeamExists проверяет, существует ли команда
	TeamExists(ctx context.Context, teamName string) (bool, error)
	// CreateTeam создает команду в БД
	CreateTeam(ctx context.Context, teamName string) error
	// GetTeam получает команду с участниками (возвращает repository/team.User)
	GetTeam(ctx context.Context, teamName string) (string, []team.User, error)
	// UserExists проверяет, существует ли пользователь
	UserExists(ctx context.Context, userID string) (bool, error)
	// CreateUser создает нового пользователя
	CreateUser(ctx context.Context, userID, username, teamName string, isActive bool) error
	// UpdateUser обновляет существующего пользователя
	UpdateUser(ctx context.Context, userID, username, teamName string, isActive bool) error
	// BeginTx начинает транзакцию (возвращает repository/team.Tx)
	BeginTx(ctx context.Context) (team.Tx, error)
}
