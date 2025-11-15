package pullrequests

import (
	"context"
	"database/sql"

	"github.com/aabbuukkaarr8/PRService/internal/repository/users"
)

// GetUser получает пользователя по user_id (используем repository/users)
func (r *Repository) GetUser(ctx context.Context, userID string) (users.User, error) {
	var user users.User

	err := r.store.GetConn().QueryRowContext(ctx,
		"SELECT user_id, username, team_name, is_active FROM users WHERE user_id = $1",
		userID).Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return users.User{}, sql.ErrNoRows
		}
		return users.User{}, err
	}

	return user, nil
}
