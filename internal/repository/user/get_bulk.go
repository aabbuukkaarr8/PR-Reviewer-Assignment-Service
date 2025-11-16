package user

import (
	"context"

	"github.com/lib/pq"
)

func (r *Repository) GetUsersTeamNames(ctx context.Context, userIDs []string) (map[string]string, error) {
	if len(userIDs) == 0 {
		return make(map[string]string), nil
	}

	query := `
		SELECT user_id, team_name
		FROM users
		WHERE user_id = ANY($1)
	`

	rows, err := r.store.GetConn().QueryContext(ctx, query, pq.Array(userIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var userID, teamName string
		if err := rows.Scan(&userID, &teamName); err != nil {
			return nil, err
		}
		result[userID] = teamName
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

