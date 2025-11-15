package pullrequests

import "time"

type PullRequest struct {
	PullRequestID     string     `json:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            string     `json:"status"` // OPEN or MERGED
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         *time.Time `json:"createdAt,omitempty"`
	MergedAt          *time.Time `json:"mergedAt,omitempty"`
}

type PullRequestShort struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"` // OPEN or MERGED
}

// CreatePullRequestRequest - запрос на создание PR
type CreatePullRequestRequest struct {
	PullRequestID   string `json:"pull_request_id" binding:"required"`
	PullRequestName string `json:"pull_request_name" binding:"required"`
	AuthorID        string `json:"author_id" binding:"required"`
}

// CreatePullRequestResponse - ответ при успешном создании PR
type CreatePullRequestResponse struct {
	PR PullRequest `json:"pr"`
}

// MergePullRequestRequest - запрос на merge PR
type MergePullRequestRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
}

// MergePullRequestResponse - ответ при успешном merge PR
type MergePullRequestResponse struct {
	PR PullRequest `json:"pr"`
}

// ReassignReviewerRequest - запрос на переназначение ревьювера
type ReassignReviewerRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
	OldUserID     string `json:"old_reviewer_id" binding:"required"`
}

// ReassignReviewerResponse - ответ при успешном переназначении
type ReassignReviewerResponse struct {
	PR         PullRequest `json:"pr"`
	ReplacedBy string      `json:"replaced_by"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail детали ошибки
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
