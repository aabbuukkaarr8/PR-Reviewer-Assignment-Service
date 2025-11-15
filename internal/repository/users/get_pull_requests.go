package users

import (
	"context"
)

// GetUserPullRequests получает PR'ы, где пользователь назначен ревьювером
func (r *Repository) GetUserPullRequests(ctx context.Context, userID string) ([]PullRequestShort, error) {
	rows, err := r.store.GetConn().QueryContext(ctx,
		`SELECT pull_request_id, pull_request_name, author_id, status 
		 FROM pullrequests 
		 WHERE $1 = ANY(assigned_reviewers)`,
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []PullRequestShort
	for rows.Next() {
		var pr PullRequestShort
		if err := rows.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status); err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return prs, nil
}
