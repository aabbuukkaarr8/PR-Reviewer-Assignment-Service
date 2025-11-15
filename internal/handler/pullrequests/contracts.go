package pullrequests

import (
	"context"

	"github.com/aabbuukkaarr8/PRService/internal/service/pullrequests"
)

// PullRequestService интерфейс для работы с pull requests
type PullRequestService interface {
	CreatePullRequest(ctx context.Context, pullRequestID, pullRequestName, authorID string) (pullrequests.PullRequest, error)
	MergePullRequest(ctx context.Context, pullRequestID string) (pullrequests.PullRequest, error)
	ReassignReviewer(ctx context.Context, pullRequestID, oldUserID string) (pullrequests.PullRequest, string, error)
}
