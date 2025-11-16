package user

import (
	"context"

	"github.com/aabbuukkaarr8/PRService/internal/repository/user"
)

type Repo interface {
	GetUser(ctx context.Context, userID string) (user.User, error)
	UpdateUserIsActive(ctx context.Context, userID string, isActive bool) error
	GetUserPullRequests(ctx context.Context, userID string) ([]user.PullRequestShort, error)
}
