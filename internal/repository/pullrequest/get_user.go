package pullrequest

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aabbuukkaarr8/PRService/internal/repository/user"
)

func (r *Repository) GetUser(ctx context.Context, userID string) (user.User, error) {
	var user user.User

	err := r.store.GetConn().QueryRowContext(ctx,
		"SELECT user_id, username, team_name, is_active FROM users WHERE user_id = $1",
		userID).Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, sql.ErrNoRows
		}
		return user, err
	}

	return user, nil
}
