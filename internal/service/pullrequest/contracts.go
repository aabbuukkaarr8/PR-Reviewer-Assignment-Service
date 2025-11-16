package pullrequest

import (
	"context"

	prrepo "github.com/aabbuukkaarr8/PRService/internal/repository/pullrequest"
	"github.com/aabbuukkaarr8/PRService/internal/repository/user"
)

// Repo интерфейс для работы с pull requests в repository слое
type Repo interface {
	// PRExists проверяет, существует ли PR
	PRExists(ctx context.Context, pullRequestID string) (bool, error)
	// GetUser получает пользователя по user_id
	GetUser(ctx context.Context, userID string) (user.User, error)
	// GetActiveTeamMembers получает активных участников команды (исключая указанного пользователя)
	GetActiveTeamMembers(ctx context.Context, teamName string, excludeUserID string) ([]user.User, error)
	// CreatePullRequest создает PR в БД
	CreatePullRequest(ctx context.Context, request *prrepo.CreatePullRequest) (prrepo.PullRequest, error)
	// GetPullRequest получает PR по ID
	GetPullRequest(ctx context.Context, pullRequestID string) (prrepo.PullRequest, error)
	// MergePullRequest помечает PR как MERGED (идемпотентная операция)
	MergePullRequest(ctx context.Context, pullRequestID string) (prrepo.PullRequest, error)
	// UpdatePullRequestReviewers обновляет список ревьюверов PR
	UpdatePullRequestReviewers(ctx context.Context, pullRequestID string, assignedReviewers []string) (prrepo.PullRequest, error)
}
