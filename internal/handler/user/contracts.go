package user

import (
	"context"

	usersrv "github.com/aabbuukkaarr8/PRService/internal/service/user"
)

type ServiceUser interface {
	SetIsActive(ctx context.Context, userID string, isActive bool) (usersrv.User, error)
	GetReview(ctx context.Context, userID string) ([]usersrv.PullRequestShort, error)
}
