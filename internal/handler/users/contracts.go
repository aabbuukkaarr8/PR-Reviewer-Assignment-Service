package users

import (
	"context"

	"github.com/aabbuukkaarr8/PRService/internal/service/users"
)

// UserService интерфейс для работы с пользователями
type UserService interface {
	SetIsActive(ctx context.Context, userID string, isActive bool) (users.User, error)
	GetReview(ctx context.Context, userID string) ([]users.PullRequestShort, error)
}
