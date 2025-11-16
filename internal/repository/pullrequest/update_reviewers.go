package pullrequest

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

func (r *Repository) UpdatePullRequestReviewers(ctx context.Context, pullRequestID string, assignedReviewers []string) (PullRequest, error) {
	_, err := r.store.GetConn().ExecContext(ctx,
		`UPDATE pullrequests SET assigned_reviewers = $1 WHERE pull_request_id = $2`,
		pq.Array(assignedReviewers), pullRequestID)
	if err != nil {
		return PullRequest{}, err
	}

	var pr PullRequest
	var reviewers pq.StringArray

	err = r.store.GetConn().QueryRowContext(ctx,
		`SELECT pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at, merged_at 
		 FROM pullrequests WHERE pull_request_id = $1`,
		pullRequestID).Scan(
		&pr.PullRequestID,
		&pr.PullRequestName,
		&pr.AuthorID,
		&pr.Status,
		&reviewers,
		&pr.CreatedAt,
		&pr.MergedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return PullRequest{}, sql.ErrNoRows
		}
		return PullRequest{}, err
	}

	pr.AssignedReviewers = []string(reviewers)

	return pr, nil
}
