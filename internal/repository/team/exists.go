package team

import "context"

// TeamExists проверяет, существует ли команда
func (r *Repository) TeamExists(ctx context.Context, teamName string) (bool, error) {
	var exists bool
	err := r.store.GetConn().QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)",
		teamName).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// UserExists проверяет, существует ли пользователь
func (r *Repository) UserExists(ctx context.Context, userID string) (bool, error) {
	var exists bool
	err := r.store.GetConn().QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)",
		userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
