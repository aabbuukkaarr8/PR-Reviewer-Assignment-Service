package user

import "context"

func (r *Repository) UpdateUserIsActive(ctx context.Context, userID string, isActive bool) error {
	_, err := r.store.GetConn().ExecContext(ctx,
		"UPDATE users SET is_active = $1 WHERE user_id = $2",
		isActive, userID)
	return err
}
