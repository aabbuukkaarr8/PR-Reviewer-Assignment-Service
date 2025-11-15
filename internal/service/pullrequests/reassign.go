package pullrequests

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aabbuukkaarr8/PRService/internal/repository/users"
)

var (
	ErrPRMerged    = errors.New("PR_MERGED")
	ErrNotAssigned = errors.New("NOT_ASSIGNED")
	ErrNoCandidate = errors.New("NO_CANDIDATE")
)

// ReassignReviewer переназначает ревьювера на другого из его команды
func (s *Service) ReassignReviewer(ctx context.Context, pullRequestID, oldUserID string) (PullRequest, string, error) {
	// 1. Получить PR
	repoPR, err := s.repo.GetPullRequest(ctx, pullRequestID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return PullRequest{}, "", ErrNotFound
		}
		return PullRequest{}, "", err
	}

	// 2. Проверить, что PR не MERGED
	if repoPR.Status == "MERGED" {
		return PullRequest{}, "", ErrPRMerged
	}

	// 3. Проверить, что oldUserID назначен ревьювером
	isAssigned := false
	for _, reviewer := range repoPR.AssignedReviewers {
		if reviewer == oldUserID {
			isAssigned = true
			break
		}
	}
	if !isAssigned {
		return PullRequest{}, "", ErrNotAssigned
	}

	// 4. Получить старого ревьювера (проверить что он существует)
	oldReviewer, err := s.repo.GetUser(ctx, oldUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return PullRequest{}, "", ErrNotFound
		}
		return PullRequest{}, "", err
	}

	// 5. Получить активных участников команды старого ревьювера
	// Исключаем: старого ревьювера, автора PR, других назначенных ревьюверов
	excludeList := []string{oldUserID, repoPR.AuthorID}
	excludeList = append(excludeList, repoPR.AssignedReviewers...)

	teamMembers, err := s.repo.GetActiveTeamMembers(ctx, oldReviewer.TeamName, oldUserID)
	if err != nil {
		return PullRequest{}, "", err
	}

	// 6. Фильтруем кандидатов: исключаем автора и других ревьюверов
	candidates := make([]users.User, 0)
	for _, member := range teamMembers {
		shouldExclude := false
		for _, excludeID := range excludeList {
			if member.UserID == excludeID {
				shouldExclude = true
				break
			}
		}
		if !shouldExclude {
			candidates = append(candidates, member)
		}
	}

	// 7. Если нет кандидатов → 409 NO_CANDIDATE
	if len(candidates) == 0 {
		return PullRequest{}, "", ErrNoCandidate
	}

	// 8. Выбрать первого кандидата
	newReviewerID := candidates[0].UserID

	// 9. Заменить в списке assigned_reviewers
	newReviewers := make([]string, len(repoPR.AssignedReviewers))
	for i, reviewer := range repoPR.AssignedReviewers {
		if reviewer == oldUserID {
			newReviewers[i] = newReviewerID
		} else {
			newReviewers[i] = reviewer
		}
	}

	// 10. Обновить PR в БД
	updatedPR, err := s.repo.UpdatePullRequestReviewers(ctx, pullRequestID, newReviewers)
	if err != nil {
		return PullRequest{}, "", err
	}

	// 11. Конвертировать repository.PullRequest в service.PullRequest
	pr := PullRequest{}
	pr.FillFromDB(&updatedPR)

	return pr, newReviewerID, nil
}
