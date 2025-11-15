package pullrequests

import (
	"context"

	repoPullRequests "github.com/aabbuukkaarr8/PRService/internal/repository/pullrequests"
	"github.com/aabbuukkaarr8/PRService/internal/repository/users"
)

// Repo интерфейс для работы с pull requests в repository слое
type Repo interface {
	// PRExists проверяет, существует ли PR
	PRExists(ctx context.Context, pullRequestID string) (bool, error)
	// GetUser получает пользователя по user_id
	GetUser(ctx context.Context, userID string) (users.User, error)
	// GetActiveTeamMembers получает активных участников команды (исключая указанного пользователя)
	GetActiveTeamMembers(ctx context.Context, teamName string, excludeUserID string) ([]users.User, error)
	// CreatePullRequest создает PR в БД
	CreatePullRequest(ctx context.Context, pullRequestID, pullRequestName, authorID, status string, assignedReviewers []string) (repoPullRequests.PullRequest, error)
	// GetPullRequest получает PR по ID
	GetPullRequest(ctx context.Context, pullRequestID string) (repoPullRequests.PullRequest, error)
	// MergePullRequest помечает PR как MERGED (идемпотентная операция)
	MergePullRequest(ctx context.Context, pullRequestID string) (repoPullRequests.PullRequest, error)
	// UpdatePullRequestReviewers обновляет список ревьюверов PR
	UpdatePullRequestReviewers(ctx context.Context, pullRequestID string, assignedReviewers []string) (repoPullRequests.PullRequest, error)
}
