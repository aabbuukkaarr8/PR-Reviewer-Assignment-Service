package pullrequests

import "time"

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
