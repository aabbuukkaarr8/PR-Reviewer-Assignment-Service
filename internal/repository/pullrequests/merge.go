package pullrequests

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

// MergePullRequest помечает PR как MERGED (идемпотентная операция)
func (r *Repository) MergePullRequest(ctx context.Context, pullRequestID string) (PullRequest, error) {
	now := time.Now()

	// Обновляем статус и mergedAt
	_, err := r.store.GetConn().ExecContext(ctx,
		`UPDATE pullrequests SET status = 'MERGED', "mergedAt" = $1 WHERE pull_request_id = $2`,
		now, pullRequestID)
	if err != nil {
		return PullRequest{}, err
	}

	// Получаем обновленный PR
	var pr PullRequest
	var assignedReviewers pq.StringArray

	err = r.store.GetConn().QueryRowContext(ctx,
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
