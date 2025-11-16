package user

import (
	"context"
	"database/sql"
	"errors"

	repoUsers "github.com/aabbuukkaarr8/PRService/internal/repository/user"
)

func (s *Service) GetReview(ctx context.Context, userID string) ([]PullRequestShort, error) {
	_, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	repoPRs, err := s.repo.GetUserPullRequests(ctx, userID)
	if err != nil {
		return nil, err
	}

	servicePRs := make([]PullRequestShort, len(repoPRs))
	for i, pr := range repoPRs {
		servicePRs[i] = toServicePullRequestShort(pr)
	}

	return servicePRs, nil
}

func toServicePullRequestShort(r repoUsers.PullRequestShort) PullRequestShort {
	return PullRequestShort{
		PullRequestID:   r.PullRequestID,
		PullRequestName: r.PullRequestName,
		AuthorID:        r.AuthorID,
		Status:          r.Status,
	}
}
