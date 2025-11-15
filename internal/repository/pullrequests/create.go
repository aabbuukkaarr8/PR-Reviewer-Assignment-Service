package pullrequests

import (
	"context"
	"time"

	"github.com/lib/pq"
)

// CreatePullRequest создает PR в БД
func (r *Repository) CreatePullRequest(ctx context.Context, pullRequestID, pullRequestName, authorID, status string, assignedReviewers []string) (PullRequest, error) {
	now := time.Now()

	_, err := r.store.GetConn().ExecContext(ctx,
		`INSERT INTO pullrequests (pull_request_id, pull_request_name, author_id, status, assigned_reviewers, "createdAt") 
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		pullRequestID, pullRequestName, authorID, status, pq.Array(assignedReviewers), now)
	if err != nil {
		return PullRequest{}, err
	}

	// Возвращаем созданный PR
	return PullRequest{
		PullRequestID:     pullRequestID,
		PullRequestName:   pullRequestName,
		AuthorID:          authorID,
		Status:            status,
		AssignedReviewers: assignedReviewers,
		CreatedAt:         &now,
		MergedAt:          nil,
	}, nil
}
