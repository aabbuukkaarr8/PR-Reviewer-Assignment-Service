package pullrequest

import (
	"context"

	prsrv "github.com/aabbuukkaarr8/PRService/internal/service/pullrequest"
)

type ServicePR interface {
	CreatePullRequest(ctx context.Context, request prsrv.CreatePullRequest) (prsrv.PullRequest, error)
	MergePullRequest(ctx context.Context, pullRequestID string) (prsrv.PullRequest, error)
	ReassignReviewer(ctx context.Context, pullRequestID, oldUserID string) (prsrv.PullRequest, string, error)
	GetStats(ctx context.Context) (prsrv.Stats, error)
	BulkDeactivateTeamUsers(ctx context.Context, teamName string) (prsrv.BulkDeactivateResult, error)
}
