package pullrequest

import (
	"time"

	"github.com/aabbuukkaarr8/PRService/internal/api/models"
)

// PullRequest структура PR для repository слоя
type PullRequest struct {
	PullRequestID     string
	PullRequestName   string
	AuthorID          string
	Status            string // OPEN or MERGED
	AssignedReviewers []string
	CreatedAt         *time.Time
	MergedAt          *time.Time
}

type CreatePullRequest struct {
	AuthorId          string
	PullRequestId     string
	PullRequestName   string
	Status            models.PullRequestStatus
	AssignedReviewers []string
}
