package user

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
	repoUser, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}

	if err := s.repo.UpdateUserIsActive(ctx, userID, isActive); err != nil {
		return User{}, err
	}

	repoUser.IsActive = isActive

	user := User{}
	user.FillFromDB(&repoUser)

	return user, nil
}
