package pullrequests

import (
	"github.com/aabbuukkaarr8/PRService/internal/api/models"
)

// CreatePullRequestRequest - запрос на создание PR
type CreatePullRequestRequest struct {
	PullRequestID   string `json:"pull_request_id" binding:"required"`
	PullRequestName string `json:"pull_request_name" binding:"required"`
	AuthorID        string `json:"author_id" binding:"required"`
}

// CreatePullRequestResponse - ответ при успешном создании PR
type CreatePullRequestResponse struct {
	PR models.PullRequest `json:"pr"`
}

// MergePullRequestRequest - запрос на merge PR
type MergePullRequestRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
}

// MergePullRequestResponse - ответ при успешном merge PR
type MergePullRequestResponse struct {
	PR models.PullRequest `json:"pr"`
}

// ReassignReviewerRequest - запрос на переназначение ревьювера
type ReassignReviewerRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
	OldUserID     string `json:"old_reviewer_id" binding:"required"`
}

// ReassignReviewerResponse - ответ при успешном переназначении
type ReassignReviewerResponse struct {
	PR         models.PullRequest `json:"pr"`
	ReplacedBy string             `json:"replaced_by"`
}

// ErrorResponse соответствует схеме ErrorResponse из OpenAPI
type ErrorResponse = models.ErrorResponse
