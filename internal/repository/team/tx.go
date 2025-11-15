package team

// Tx интерфейс для работы с транзакцией (определен в repository слое)
type Tx interface {
	// CreateTeam создает команду в транзакции
	CreateTeam(teamName string) error
	// CreateUser создает пользователя в транзакции
	CreateUser(userID, username, teamName string, isActive bool) error
	// UpdateUser обновляет пользователя в транзакции
	UpdateUser(userID, username, teamName string, isActive bool) error
	// Commit коммитит транзакцию
	Commit() error
	// Rollback откатывает транзакцию
	Rollback() error
}
