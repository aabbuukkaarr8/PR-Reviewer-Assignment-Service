package pullrequest

import (
	"context"
	"errors"
	"math/rand"
	"time"

	prrepo "github.com/aabbuukkaarr8/PRService/internal/repository/pullrequest"
)

type BulkDeactivateResult struct {
	DeactivatedUserIDs []string
	ReassignedPRs      []ReassignedPR
}

type ReassignedPR struct {
	PullRequestID string
	OldReviewerID string
	NewReviewerID string
}

func (s *Service) BulkDeactivateTeamUsers(ctx context.Context, teamName string) (BulkDeactivateResult, error) {
	deactivatedUserIDs, err := s.repo.BulkDeactivateTeamUsers(ctx, teamName)
	if err != nil {
		return BulkDeactivateResult{}, err
	}

	if len(deactivatedUserIDs) == 0 {
		return BulkDeactivateResult{
			DeactivatedUserIDs: []string{},
			ReassignedPRs:      []ReassignedPR{},
		}, nil
	}

	openPRs, err := s.repo.GetOpenPRsByReviewers(ctx, deactivatedUserIDs)
	if err != nil {
		return BulkDeactivateResult{}, err
	}

	if len(openPRs) == 0 {
		return BulkDeactivateResult{
			DeactivatedUserIDs: deactivatedUserIDs,
			ReassignedPRs:      []ReassignedPR{},
		}, nil
	}

	teamMembers, err := s.repo.GetActiveTeamMembers(ctx, teamName, "")
	if err != nil {
		return BulkDeactivateResult{}, err
	}

	activeUserIDs := make(map[string]bool)
	for _, member := range teamMembers {
		activeUserIDs[member.UserID] = true
	}

	deactivatedSet := make(map[string]bool)
	for _, userID := range deactivatedUserIDs {
		deactivatedSet[userID] = true
	}

	var reassignedPRs []ReassignedPR
	var prUpdates []prrepo.PRReviewerUpdate

	for _, pr := range openPRs {
		if pr.AuthorTeamName != teamName {
			continue
		}

		needsReassignment := false
		newReviewers := make([]string, 0, len(pr.AssignedReviewers))
		replacedReviewers := make(map[string]string)

		for _, reviewerID := range pr.AssignedReviewers {
			if deactivatedSet[reviewerID] {
				needsReassignment = true
				candidate, err := s.findReplacementForDeactivated(ctx, teamName, reviewerID, pr.AuthorID, pr.AssignedReviewers, activeUserIDs, replacedReviewers)
				if err != nil {
					if errors.Is(err, ErrNoCandidate) {
						continue
					}
					return BulkDeactivateResult{}, err
				}
				newReviewers = append(newReviewers, candidate)
				replacedReviewers[reviewerID] = candidate
				reassignedPRs = append(reassignedPRs, ReassignedPR{
					PullRequestID: pr.PullRequestID,
					OldReviewerID: reviewerID,
					NewReviewerID: candidate,
				})
			} else {
				newReviewers = append(newReviewers, reviewerID)
			}
		}

		if needsReassignment && len(newReviewers) > 0 {
			prUpdates = append(prUpdates, prrepo.PRReviewerUpdate{
				PullRequestID:     pr.PullRequestID,
				AssignedReviewers: newReviewers,
			})
		}
	}

	if len(prUpdates) > 0 {
		err := s.repo.BulkUpdatePullRequestReviewers(ctx, prUpdates)
		if err != nil {
			return BulkDeactivateResult{}, err
		}
	}

	return BulkDeactivateResult{
		DeactivatedUserIDs: deactivatedUserIDs,
		ReassignedPRs:      reassignedPRs,
	}, nil
}

func (s *Service) findReplacementForDeactivated(ctx context.Context, teamName, oldUserID, authorID string, currentReviewers []string, activeUserIDs map[string]bool, alreadyReplaced map[string]string) (string, error) {
	excludeList := make(map[string]bool)
	excludeList[oldUserID] = true
	excludeList[authorID] = true
	for _, reviewer := range currentReviewers {
		excludeList[reviewer] = true
	}
	for _, replaced := range alreadyReplaced {
		excludeList[replaced] = true
	}

	candidates := make([]string, 0)
	for userID := range activeUserIDs {
		if !excludeList[userID] {
			candidates = append(candidates, userID)
		}
	}

	if len(candidates) == 0 {
		return "", ErrNoCandidate
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	selected := candidates[r.Intn(len(candidates))]

	return selected, nil
}
