package pullrequest

import (
	"time"

	prsrv "github.com/aabbuukkaarr8/PRService/internal/service/pullrequest"
)

type PullRequest struct {
	PullRequestID     string     `json:"pull_request_id" binding:"required"`
	PullRequestName   string     `json:"pull_request_name" binding:"required"`
	AuthorID          string     `json:"author_id" binding:"required"`
	Status            string     `json:"status" binding:"required,oneof=OPEN MERGED"`
	AssignedReviewers []string   `json:"assigned_reviewers" binding:"max=2,dive,required"`
	CreatedAt         *time.Time `json:"created_at,omitempty"`
	MergedAt          *time.Time `json:"merged_at,omitempty"`
}

type CreatePullRequestResponse struct {
	PR PullRequest `json:"pr"`
}

type MergeRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
}

type MergePullRequestResponse struct {
	PR PullRequest `json:"pr"`
}

type ReassignReviewerRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
	OldUserID     string `json:"old_user_id" binding:"required"`
}

type ReassignReviewerResponse struct {
	PR         PullRequest `json:"pr"`
	ReplacedBy string      `json:"replaced_by"`
}

func toHandlerPullRequest(s prsrv.PullRequest) PullRequest {
	return PullRequest{
		PullRequestID:     s.PullRequestID,
		PullRequestName:   s.PullRequestName,
		AuthorID:          s.AuthorID,
		Status:            s.Status,
		AssignedReviewers: s.AssignedReviewers,
		CreatedAt:         s.CreatedAt,
		MergedAt:          s.MergedAt,
	}
}
