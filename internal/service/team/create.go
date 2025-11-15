package team

import (
	"context"
	"errors"
)

var (
	ErrTeamExists = errors.New("TEAM_EXISTS")
)

// CreateTeam создает команду и пользователей
func (s *Service) CreateTeam(ctx context.Context, team Team) (Team, error) {
	// 1. Проверить, существует ли команда
	exists, err := s.repo.TeamExists(ctx, team.TeamName)
	if err != nil {
		return Team{}, err
	}
	if exists {
		return Team{}, ErrTeamExists
	}

	// 2. Начать транзакцию
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return Team{}, err
	}
	defer tx.Rollback() // Откатим, если что-то пойдет не так

	// 3. Создать команду
	if err := tx.CreateTeam(team.TeamName); err != nil {
		return Team{}, err
	}

	// 4. Для каждого участника: создать или обновить
	for _, member := range team.Members {
		// Проверяем, существует ли пользователь
		userExists, err := s.repo.UserExists(ctx, member.UserID)
		if err != nil {
			return Team{}, err
		}

		if userExists {
			// Обновляем существующего пользователя
			if err := tx.UpdateUser(member.UserID, member.Username, team.TeamName, member.IsActive); err != nil {
				return Team{}, err
			}
		} else {
			// Создаем нового пользователя
			if err := tx.CreateUser(member.UserID, member.Username, team.TeamName, member.IsActive); err != nil {
				return Team{}, err
			}
		}
	}

	// 5. Закоммитить транзакцию
	if err := tx.Commit(); err != nil {
		return Team{}, err
	}

	// 6. Вернуть созданную команду
	return team, nil
}
