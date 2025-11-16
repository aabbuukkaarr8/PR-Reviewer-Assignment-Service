package pullrequest

import (
	"time"

	"github.com/aabbuukkaarr8/PRService/internal/api/models"
)

type PullRequest struct {
	PullRequestID     string
	PullRequestName   string
	AuthorID          string
	Status            string
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
