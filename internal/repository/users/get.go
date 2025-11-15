package users

import (
	"context"
	"database/sql"
)

// GetUser получает пользователя по user_id
func (r *Repository) GetUser(ctx context.Context, userID string) (User, error) {
	var user User

	err := r.store.GetConn().QueryRowContext(ctx,
		"SELECT user_id, username, team_name, is_active FROM users WHERE user_id = $1",
		userID).Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, sql.ErrNoRows
		}
		return User{}, err
	}

	return user, nil
}
