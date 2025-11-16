package pullrequest

import (
	"context"
	"database/sql"
	"errors"
)

// MergePullRequest помечает PR как MERGED (идемпотентная операция)
func (s *Service) MergePullRequest(ctx context.Context, pullRequestID string) (PullRequest, error) {
	repoPR, err := s.repo.GetPullRequest(ctx, pullRequestID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return PullRequest{}, ErrNotFound
		}
		return PullRequest{}, err
	}

	if repoPR.Status == "MERGED" {
		pr := PullRequest{}
		pr.FillFromDB(&repoPR)
		return pr, nil
	}

	mergedPR, err := s.repo.MergePullRequest(ctx, pullRequestID)
	if err != nil {
		return PullRequest{}, err
	}

	pr := PullRequest{}
	pr.FillFromDB(&mergedPR)

	return pr, nil
}
