package pullrequests

import (
	"context"
	"database/sql"
	"errors"
)

// MergePullRequest помечает PR как MERGED (идемпотентная операция)
func (s *Service) MergePullRequest(ctx context.Context, pullRequestID string) (PullRequest, error) {
	// 1. Получить PR
	repoPR, err := s.repo.GetPullRequest(ctx, pullRequestID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return PullRequest{}, ErrNotFound
		}
		return PullRequest{}, err
	}

	// 2. Если PR уже MERGED, просто вернуть его (идемпотентность)
	if repoPR.Status == "MERGED" {
		pr := PullRequest{}
		pr.FillFromDB(&repoPR)
		return pr, nil
	}

	// 3. Пометить PR как MERGED
	mergedPR, err := s.repo.MergePullRequest(ctx, pullRequestID)
	if err != nil {
		return PullRequest{}, err
	}

	// 4. Конвертировать repository.PullRequest в service.PullRequest
	pr := PullRequest{}
	pr.FillFromDB(&mergedPR)

	return pr, nil
}
