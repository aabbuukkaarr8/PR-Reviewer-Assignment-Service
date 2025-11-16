package pullrequest

import (
	"time"

	prrepo "github.com/aabbuukkaarr8/PRService/internal/repository/pullrequest"
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
	AuthorId        string
	PullRequestId   string
	PullRequestName string
}

func (m *PullRequest) FillFromDB(dbp *prrepo.PullRequest) {
	m.PullRequestID = dbp.PullRequestID
	m.PullRequestName = dbp.PullRequestName
	m.AuthorID = dbp.AuthorID
	m.Status = dbp.Status
	m.AssignedReviewers = dbp.AssignedReviewers
	m.CreatedAt = dbp.CreatedAt
	m.MergedAt = dbp.MergedAt
}

func (m *CreatePullRequest) ToDB() *prrepo.CreatePullRequest {
	return &prrepo.CreatePullRequest{
		AuthorId:        m.AuthorId,
		PullRequestId:   m.PullRequestId,
		PullRequestName: m.PullRequestName,
	}
}
