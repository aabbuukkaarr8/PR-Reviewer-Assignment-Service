package users

import "context"

// UpdateUserIsActive обновляет только флаг is_active пользователя
func (r *Repository) UpdateUserIsActive(ctx context.Context, userID string, isActive bool) error {
	_, err := r.store.GetConn().ExecContext(ctx,
		"UPDATE users SET is_active = $1 WHERE user_id = $2",
		isActive, userID)
	return err
}
