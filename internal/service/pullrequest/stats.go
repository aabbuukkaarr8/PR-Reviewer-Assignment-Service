package pullrequest

import (
	"context"
)

type ReviewerStats struct {
	UserID          string
	Username        string
	TeamName        string
	AssignmentsCount int
}

type PRStats struct {
	TotalPRs  int
	OpenPRs   int
	MergedPRs int
}

type Stats struct {
	PRStats        PRStats
	ReviewerStats  []ReviewerStats
}

func (s *Service) GetStats(ctx context.Context) (Stats, error) {
	prStats, err := s.repo.GetPRStats(ctx)
	if err != nil {
		return Stats{}, err
	}

	reviewerStats, err := s.repo.GetReviewerStats(ctx)
	if err != nil {
		return Stats{}, err
	}

	serviceReviewerStats := make([]ReviewerStats, len(reviewerStats))
	for i, rs := range reviewerStats {
		serviceReviewerStats[i] = ReviewerStats{
			UserID:          rs.UserID,
			Username:        rs.Username,
			TeamName:        rs.TeamName,
			AssignmentsCount: rs.AssignmentsCount,
		}
	}

	return Stats{
		PRStats: PRStats{
			TotalPRs:  prStats.TotalPRs,
			OpenPRs:   prStats.OpenPRs,
			MergedPRs: prStats.MergedPRs,
		},
		ReviewerStats: serviceReviewerStats,
	}, nil
}

