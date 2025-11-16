package team

import (
	"context"
	"database/sql"
)

func (r *Repository) GetTeam(ctx context.Context, teamName string) (string, []User, error) {
	var exists bool
	err := r.store.GetConn().QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)",
		teamName).Scan(&exists)
	if err != nil {
		return "", nil, err
	}
	if !exists {
		return "", nil, sql.ErrNoRows
	}

	rows, err := r.store.GetConn().QueryContext(ctx,
		"SELECT user_id, username, is_active FROM users WHERE team_name = $1",
		teamName)
	if err != nil {
		return "", nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.UserID, &u.Username, &u.IsActive); err != nil {
			return "", nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return "", nil, err
	}

	return teamName, users, nil
}
