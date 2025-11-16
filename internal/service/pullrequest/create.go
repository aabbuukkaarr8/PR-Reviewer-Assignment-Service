package pullrequest

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"time"

	"github.com/aabbuukkaarr8/PRService/internal/api/models"
	"github.com/aabbuukkaarr8/PRService/internal/repository/user"
)

var (
	ErrPRExists = errors.New("PR_EXISTS")
	ErrNotFound = errors.New("NOT_FOUND")
)

// CreatePullRequest создает PR и автоматически назначает до 2 ревьюверов из команды автора
func (s *Service) CreatePullRequest(ctx context.Context, req CreatePullRequest) (PullRequest, error) {
	exists, err := s.repo.PRExists(ctx, req.PullRequestId)
	if err != nil {
		return PullRequest{}, err
	}
	if exists {
		return PullRequest{}, ErrPRExists
	}

	author, err := s.repo.GetUser(ctx, req.AuthorId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return PullRequest{}, ErrNotFound
		}
		return PullRequest{}, err
	}

	teamMembers, err := s.repo.GetActiveTeamMembers(ctx, author.TeamName, req.AuthorId)
	if err != nil {
		return PullRequest{}, err
	}

	// 4. Выбрать до 2 ревьюверов случайным образом
	assignedReviewers := selectRandomReviewers(teamMembers, 2)

	reqToDB := req.ToDB()
	reqToDB.Status = models.PullRequestStatusOPEN
	reqToDB.AssignedReviewers = assignedReviewers

	repoPR, err := s.repo.CreatePullRequest(ctx, reqToDB)
	if err != nil {
		return PullRequest{}, err
	}

	pr := PullRequest{}
	pr.FillFromDB(&repoPR)

	return pr, nil
}

func selectRandomReviewers(members []user.User, maxReviewers int) []string {
	if len(members) == 0 {
		return []string{}
	}

	// Если участников меньше или равно maxReviewers, возвращаем всех
	if len(members) <= maxReviewers {
		result := make([]string, len(members))
		for i, member := range members {
			result[i] = member.UserID
		}
		return result
	}

	shuffled := make([]user.User, len(members))
	copy(shuffled, members)

	// Перемешиваем случайным образом
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	result := make([]string, maxReviewers)
	for i := 0; i < maxReviewers; i++ {
		result[i] = shuffled[i].UserID
	}

	return result
}
