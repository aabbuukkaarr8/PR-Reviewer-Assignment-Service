package pullrequests

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

// GetPullRequest получает PR по ID
func (r *Repository) GetPullRequest(ctx context.Context, pullRequestID string) (PullRequest, error) {
	var pr PullRequest
	var assignedReviewers pq.StringArray

	err := r.store.GetConn().QueryRowContext(ctx,
		`SELECT pull_request_id, pull_request_name, author_id, status, assigned_reviewers, "createdAt", "mergedAt" 
		 FROM pullrequests WHERE pull_request_id = $1`,
		pullRequestID).Scan(
		&pr.PullRequestID,
		&pr.PullRequestName,
		&pr.AuthorID,
		&pr.Status,
		&assignedReviewers,
		&pr.CreatedAt,
		&pr.MergedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return PullRequest{}, sql.ErrNoRows
		}
		return PullRequest{}, err
	}

	// Конвертируем pq.StringArray в []string
	pr.AssignedReviewers = []string(assignedReviewers)

	return pr, nil
}
