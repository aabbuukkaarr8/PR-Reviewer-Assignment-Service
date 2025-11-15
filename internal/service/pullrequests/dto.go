package pullrequests

import (
	"time"

	"github.com/aabbuukkaarr8/PRService/internal/repository/pullrequests"
)

// PullRequest структура PR для service слоя
type PullRequest struct {
	PullRequestID     string
	PullRequestName   string
	AuthorID          string
	Status            string // OPEN or MERGED
	AssignedReviewers []string
	CreatedAt         *time.Time
	MergedAt          *time.Time
}

// FillFromDB конвертирует repository.PullRequest в service.PullRequest
func (m *PullRequest) FillFromDB(dbp *pullrequests.PullRequest) {
	m.PullRequestID = dbp.PullRequestID
	m.PullRequestName = dbp.PullRequestName
	m.AuthorID = dbp.AuthorID
	m.Status = dbp.Status
	m.AssignedReviewers = dbp.AssignedReviewers
	m.CreatedAt = dbp.CreatedAt
	m.MergedAt = dbp.MergedAt
}
