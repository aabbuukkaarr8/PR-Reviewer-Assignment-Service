package users

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrUserNotFound = errors.New("NOT_FOUND")
)

// SetIsActive устанавливает флаг активности пользователя
func (s *Service) SetIsActive(ctx context.Context, userID string, isActive bool) (User, error) {
	// Получаем пользователя из repository (RepoUser)
	repoUser, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}

	// Обновляем флаг is_active в БД
	if err := s.repo.UpdateUserIsActive(ctx, userID, isActive); err != nil {
		return User{}, err
	}

	// Обновляем флаг в структуре repoUser (так как мы уже обновили в БД)
	repoUser.IsActive = isActive

	// Конвертируем RepoUser в service.User
	user := User{}
	user.FillFromDB(&repoUser)

	return user, nil
}
