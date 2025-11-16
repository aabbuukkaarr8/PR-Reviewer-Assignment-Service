package pullrequest

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

func (r *Repository) MergePullRequest(ctx context.Context, pullRequestID string) (PullRequest, error) {
	now := time.Now()

	_, err := r.store.GetConn().ExecContext(ctx,
		`UPDATE pullrequests SET status = 'MERGED', merged_at = $1 WHERE pull_request_id = $2`,
		now, pullRequestID)
	if err != nil {
		return PullRequest{}, err
	}

	var pr PullRequest
	var assignedReviewers pq.StringArray

	err = r.store.GetConn().QueryRowContext(ctx,
		`SELECT pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at, merged_at 
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

	pr.AssignedReviewers = []string(assignedReviewers)

	return pr, nil
}
