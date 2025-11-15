package users

import (
	"context"
	"database/sql"
	"errors"

	repoUsers "github.com/aabbuukkaarr8/PRService/internal/repository/users"
)

// GetReview получает PR'ы, где пользователь назначен ревьювером
func (s *Service) GetReview(ctx context.Context, userID string) ([]PullRequestShort, error) {
	// 1. Проверить, что пользователь существует
	_, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// 2. Получить PR'ы пользователя из repository
	repoPRs, err := s.repo.GetUserPullRequests(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 3. Конвертировать repository.PullRequestShort в service.PullRequestShort
	servicePRs := make([]PullRequestShort, len(repoPRs))
	for i, pr := range repoPRs {
		servicePRs[i] = toServicePullRequestShort(pr)
	}

	return servicePRs, nil
}

// toServicePullRequestShort конвертирует repository.PullRequestShort в service.PullRequestShort
func toServicePullRequestShort(r repoUsers.PullRequestShort) PullRequestShort {
	return PullRequestShort{
		PullRequestID:   r.PullRequestID,
		PullRequestName: r.PullRequestName,
		AuthorID:        r.AuthorID,
		Status:          r.Status,
	}
}
