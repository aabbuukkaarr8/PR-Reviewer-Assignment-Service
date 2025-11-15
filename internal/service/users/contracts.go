package users

import (
	"context"

	"github.com/aabbuukkaarr8/PRService/internal/repository/users"
)

type Repo interface {
	// GetUser получает пользователя по user_id
	GetUser(ctx context.Context, userID string) (users.User, error)
	// UpdateUserIsActive обновляет только флаг is_active пользователя
	UpdateUserIsActive(ctx context.Context, userID string, isActive bool) error
	// GetUserPullRequests получает PR'ы, где пользователь назначен ревьювером
	GetUserPullRequests(ctx context.Context, userID string) ([]users.PullRequestShort, error)
}
