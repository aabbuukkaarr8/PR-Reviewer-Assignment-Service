package pullrequest

import (
	"context"

	"github.com/aabbuukkaarr8/PRService/internal/repository/user"
)

// GetActiveTeamMembers получает активных участников команды (исключая указанного пользователя)
func (r *Repository) GetActiveTeamMembers(ctx context.Context, teamName string, excludeUserID string) ([]user.User, error) {
	rows, err := r.store.GetConn().QueryContext(ctx,
		"SELECT user_id, username, team_name, is_active FROM users WHERE team_name = $1 AND is_active = TRUE AND user_id != $2",
		teamName, excludeUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []user.User
	for rows.Next() {
		var u user.User
		if err := rows.Scan(&u.UserID, &u.Username, &u.TeamName, &u.IsActive); err != nil {
			return nil, err
		}
		members = append(members, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}
