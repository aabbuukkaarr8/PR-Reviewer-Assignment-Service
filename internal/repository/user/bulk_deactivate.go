package user

import (
	"context"
)

func (r *Repository) BulkDeactivateTeamUsers(ctx context.Context, teamName string) ([]string, error) {
	query := `
		UPDATE users
		SET is_active = FALSE
		WHERE team_name = $1 AND is_active = TRUE
		RETURNING user_id
	`

	rows, err := r.store.GetConn().QueryContext(ctx, query, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deactivatedUserIDs []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		deactivatedUserIDs = append(deactivatedUserIDs, userID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return deactivatedUserIDs, nil
}

