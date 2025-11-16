package pullrequest

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aabbuukkaarr8/PRService/internal/repository/user"
)

var (
	ErrPRMerged    = errors.New("PR_MERGED")
	ErrNotAssigned = errors.New("NOT_ASSIGNED")
	ErrNoCandidate = errors.New("NO_CANDIDATE")
)

func (s *Service) ReassignReviewer(ctx context.Context, pullRequestID, oldUserID string) (PullRequest, string, error) {
	repoPR, err := s.repo.GetPullRequest(ctx, pullRequestID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return PullRequest{}, "", ErrNotFound
		}
		return PullRequest{}, "", err
	}

	if repoPR.Status == "MERGED" {
		return PullRequest{}, "", ErrPRMerged
	}

	if !isReviewerAssigned(repoPR.AssignedReviewers, oldUserID) {
		return PullRequest{}, "", ErrNotAssigned
	}

	oldReviewer, err := s.repo.GetUser(ctx, oldUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return PullRequest{}, "", ErrNotFound
		}
		return PullRequest{}, "", err
	}

	candidates, err := s.findReplacementCandidates(ctx, oldReviewer.TeamName, oldUserID, repoPR.AuthorID, repoPR.AssignedReviewers)
	if err != nil {
		return PullRequest{}, "", err
	}

	if len(candidates) == 0 {
		return PullRequest{}, "", ErrNoCandidate
	}

	newReviewerID, err := selectReplacementReviewer(candidates)
	if err != nil {
		return PullRequest{}, "", err
	}

	newReviewers := replaceReviewerInList(repoPR.AssignedReviewers, oldUserID, newReviewerID)

	updatedPR, err := s.repo.UpdatePullRequestReviewers(ctx, pullRequestID, newReviewers)
	if err != nil {
		return PullRequest{}, "", err
	}

	pr := PullRequest{}
	pr.FillFromDB(&updatedPR)

	return pr, newReviewerID, nil
}

// isReviewerAssigned проверяет, назначен ли пользователь ревьювером PR
func isReviewerAssigned(assignedReviewers []string, userID string) bool {
	for _, reviewer := range assignedReviewers {
		if reviewer == userID {
			return true
		}
	}
	return false
}

// findReplacementCandidates находит кандидатов для замены ревьювера
// Исключает: старого ревьювера, автора PR, других назначенных ревьюверов
func (s *Service) findReplacementCandidates(ctx context.Context, teamName, oldUserID, authorID string, assignedReviewers []string) ([]user.User, error) {
	teamMembers, err := s.repo.GetActiveTeamMembers(ctx, teamName, oldUserID)
	if err != nil {
		return nil, err
	}

	excludeList := make(map[string]bool)
	excludeList[oldUserID] = true
	excludeList[authorID] = true
	for _, reviewer := range assignedReviewers {
		excludeList[reviewer] = true
	}

	candidates := make([]user.User, 0)
	for _, member := range teamMembers {
		if !excludeList[member.UserID] {
			candidates = append(candidates, member)
		}
	}

	return candidates, nil
}

// selectReplacementReviewer выбирает случайного кандидата для замены
func selectReplacementReviewer(candidates []user.User) (string, error) {
	if len(candidates) == 0 {
		return "", ErrNoCandidate
	}

	selected := selectRandomReviewers(candidates, 1)
	if len(selected) == 0 {
		return "", ErrNoCandidate
	}

	return selected[0], nil
}

// replaceReviewerInList заменяет старого ревьювера на нового в списке
func replaceReviewerInList(assignedReviewers []string, oldUserID, newUserID string) []string {
	newReviewers := make([]string, len(assignedReviewers))
	for i, reviewer := range assignedReviewers {
		if reviewer == oldUserID {
			newReviewers[i] = newUserID
		} else {
			newReviewers[i] = reviewer
		}
	}
	return newReviewers
}
