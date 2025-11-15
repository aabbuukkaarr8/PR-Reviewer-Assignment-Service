package team

import "context"

// CreateTeam создает команду в БД (без транзакции)
func (r *Repository) CreateTeam(ctx context.Context, teamName string) error {
	_, err := r.store.GetConn().ExecContext(ctx,
		"INSERT INTO teams (team_name) VALUES ($1)",
		teamName)
	return err
}

// CreateUser создает нового пользователя (без транзакции)
func (r *Repository) CreateUser(ctx context.Context, userID, username, teamName string, isActive bool) error {
	_, err := r.store.GetConn().ExecContext(ctx,
		"INSERT INTO users (user_id, username, team_name, is_active) VALUES ($1, $2, $3, $4)",
		userID, username, teamName, isActive)
	return err
}
