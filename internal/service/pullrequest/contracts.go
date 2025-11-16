package pullrequest

import (
	"context"

	prrepo "github.com/aabbuukkaarr8/PRService/internal/repository/pullrequest"
	"github.com/aabbuukkaarr8/PRService/internal/repository/user"
)

type Repo interface {
	PRExists(ctx context.Context, pullRequestID string) (bool, error)
	GetUser(ctx context.Context, userID string) (user.User, error)
	GetActiveTeamMembers(ctx context.Context, teamName string, excludeUserID string) ([]user.User, error)
	CreatePullRequest(ctx context.Context, request *prrepo.CreatePullRequest) (prrepo.PullRequest, error)
	GetPullRequest(ctx context.Context, pullRequestID string) (prrepo.PullRequest, error)
	MergePullRequest(ctx context.Context, pullRequestID string) (prrepo.PullRequest, error)
	UpdatePullRequestReviewers(ctx context.Context, pullRequestID string, assignedReviewers []string) (prrepo.PullRequest, error)
	GetReviewerStats(ctx context.Context) ([]prrepo.ReviewerStats, error)
	GetPRStats(ctx context.Context) (prrepo.PRStats, error)
	GetOpenPRsByReviewers(ctx context.Context, reviewerIDs []string) ([]prrepo.OpenPRWithReviewer, error)
	BulkDeactivateTeamUsers(ctx context.Context, teamName string) ([]string, error)
	BulkUpdatePullRequestReviewers(ctx context.Context, updates []prrepo.PRReviewerUpdate) error
}
