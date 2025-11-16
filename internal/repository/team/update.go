package team

import "context"

func (r *Repository) UpdateUser(ctx context.Context, userID, username, teamName string, isActive bool) error {
	_, err := r.store.GetConn().ExecContext(ctx,
		"UPDATE users SET username = $1, team_name = $2, is_active = $3 WHERE user_id = $4",
		username, teamName, isActive, userID)
	return err
}
