package pullrequests

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrPRExists = errors.New("PR_EXISTS")
	ErrNotFound = errors.New("NOT_FOUND")
)

// CreatePullRequest создает PR и автоматически назначает до 2 ревьюверов из команды автора
func (s *Service) CreatePullRequest(ctx context.Context, pullRequestID, pullRequestName, authorID string) (PullRequest, error) {
	// 1. Проверить, существует ли PR
	exists, err := s.repo.PRExists(ctx, pullRequestID)
	if err != nil {
		return PullRequest{}, err
	}
	if exists {
		return PullRequest{}, ErrPRExists
	}

	// 2. Получить автора (проверить что он существует)
	author, err := s.repo.GetUser(ctx, authorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return PullRequest{}, ErrNotFound
		}
		return PullRequest{}, err
	}

	// 3. Получить активных участников команды автора (исключая автора)
	teamMembers, err := s.repo.GetActiveTeamMembers(ctx, author.TeamName, authorID)
	if err != nil {
		return PullRequest{}, err
	}

	// 4. Выбрать до 2 ревьюверов
	assignedReviewers := make([]string, 0, 2)
	for i, member := range teamMembers {
		if i >= 2 {
			break
		}
		assignedReviewers = append(assignedReviewers, member.UserID)
	}

	// 5. Создать PR со статусом OPEN
	now := time.Now()
	repoPR, err := s.repo.CreatePullRequest(ctx, pullRequestID, pullRequestName, authorID, "OPEN", assignedReviewers)
	if err != nil {
		return PullRequest{}, err
	}

	// 6. Конвертировать repository.PullRequest в service.PullRequest
	pr := PullRequest{}
	pr.FillFromDB(&repoPR)
	pr.CreatedAt = &now

	return pr, nil
}
