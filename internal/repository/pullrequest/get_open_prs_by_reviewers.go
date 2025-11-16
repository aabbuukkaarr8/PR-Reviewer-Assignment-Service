package pullrequest

import (
	"context"

	"github.com/lib/pq"
)

type OpenPRWithReviewer struct {
	PullRequestID     string
	PullRequestName   string
	AuthorID          string
	AssignedReviewers []string
	AuthorTeamName    string
}

func (r *Repository) GetOpenPRsByReviewers(ctx context.Context, reviewerIDs []string) ([]OpenPRWithReviewer, error) {
	if len(reviewerIDs) == 0 {
		return []OpenPRWithReviewer{}, nil
	}

	query := `
		SELECT 
			pr.pull_request_id,
			pr.pull_request_name,
			pr.author_id,
			pr.assigned_reviewers,
			u.team_name as author_team_name
		FROM pullrequests pr
		INNER JOIN users u ON pr.author_id = u.user_id
		WHERE pr.status = 'OPEN'
		  AND pr.assigned_reviewers && $1
	`

	rows, err := r.store.GetConn().QueryContext(ctx, query, pq.Array(reviewerIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []OpenPRWithReviewer
	for rows.Next() {
		var pr OpenPRWithReviewer
		var reviewers pq.StringArray

		if err := rows.Scan(
			&pr.PullRequestID,
			&pr.PullRequestName,
			&pr.AuthorID,
			&reviewers,
			&pr.AuthorTeamName,
		); err != nil {
			return nil, err
		}

		pr.AssignedReviewers = []string(reviewers)

		prs = append(prs, pr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return prs, nil
}

