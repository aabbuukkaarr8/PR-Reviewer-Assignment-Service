package pullrequest

import (
	"context"
	"time"

	"github.com/lib/pq"
)

// CreatePullRequest создает PR в БД
func (r *Repository) CreatePullRequest(
	ctx context.Context,
	req *CreatePullRequest,
) (PullRequest, error) {
	now := time.Now()

	_, err := r.store.GetConn().ExecContext(ctx,
		`INSERT INTO pullrequests (pull_request_id, pull_request_name, author_id, status, assigned_reviewers, 
"created_at") 
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		req.PullRequestId, req.PullRequestName, req.AuthorId, string(req.Status), pq.Array(req.AssignedReviewers), now)
	if err != nil {
		return PullRequest{}, err
	}

	// Возвращаем созданный PR
	return PullRequest{
		PullRequestID:     req.PullRequestId,
		PullRequestName:   req.PullRequestName,
		AuthorID:          req.AuthorId,
		Status:            string(req.Status),
		AssignedReviewers: req.AssignedReviewers,
		CreatedAt:         &now,
		MergedAt:          nil,
	}, nil
}
