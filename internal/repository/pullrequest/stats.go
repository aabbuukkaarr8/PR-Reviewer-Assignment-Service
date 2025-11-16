package pullrequest

import (
	"context"
)

type ReviewerStats struct {
	UserID        string
	Username      string
	TeamName      string
	AssignmentsCount int
}

func (r *Repository) GetReviewerStats(ctx context.Context) ([]ReviewerStats, error) {
	query := `
		SELECT 
			u.user_id,
			u.username,
			u.team_name,
			COUNT(*) as assignments_count
		FROM users u
		INNER JOIN pullrequests pr ON u.user_id = ANY(pr.assigned_reviewers)
		WHERE u.is_active = TRUE
		GROUP BY u.user_id, u.username, u.team_name
		ORDER BY assignments_count DESC
	`

	rows, err := r.store.GetConn().QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []ReviewerStats
	for rows.Next() {
		var s ReviewerStats
		if err := rows.Scan(&s.UserID, &s.Username, &s.TeamName, &s.AssignmentsCount); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}

type PRStats struct {
	TotalPRs    int
	OpenPRs     int
	MergedPRs   int
}

func (r *Repository) GetPRStats(ctx context.Context) (PRStats, error) {
	var stats PRStats

	err := r.store.GetConn().QueryRowContext(ctx,
		`SELECT 
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'OPEN') as open,
			COUNT(*) FILTER (WHERE status = 'MERGED') as merged
		FROM pullrequests`).Scan(&stats.TotalPRs, &stats.OpenPRs, &stats.MergedPRs)
	if err != nil {
		return PRStats{}, err
	}

	return stats, nil
}

