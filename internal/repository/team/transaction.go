package team

import (
	"context"
	"database/sql"
)

// Transaction обертка над *sql.Tx для реализации Tx
type Transaction struct {
	tx *sql.Tx
}

// BeginTx начинает транзакцию
func (r *Repository) BeginTx(ctx context.Context) (Tx, error) {
	tx, err := r.store.GetConn().BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Transaction{tx: tx}, nil
}

// CreateTeam создает команду в транзакции
func (t *Transaction) CreateTeam(teamName string) error {
	_, err := t.tx.Exec("INSERT INTO teams (team_name) VALUES ($1)", teamName)
	return err
}

// CreateUser создает пользователя в транзакции
func (t *Transaction) CreateUser(userID, username, teamName string, isActive bool) error {
	_, err := t.tx.Exec(
		"INSERT INTO users (user_id, username, team_name, is_active) VALUES ($1, $2, $3, $4)",
		userID, username, teamName, isActive)
	return err
}

// UpdateUser обновляет пользователя в транзакции
func (t *Transaction) UpdateUser(userID, username, teamName string, isActive bool) error {
	_, err := t.tx.Exec(
		"UPDATE users SET username = $1, team_name = $2, is_active = $3 WHERE user_id = $4",
		username, teamName, isActive, userID)
	return err
}

// Commit коммитит транзакцию
func (t *Transaction) Commit() error {
	return t.tx.Commit()
}

// Rollback откатывает транзакцию
func (t *Transaction) Rollback() error {
	return t.tx.Rollback()
}
