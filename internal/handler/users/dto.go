package users

import (
	"github.com/aabbuukkaarr8/PRService/internal/api/models"
)

// SetIsActiveRequest - запрос на установку флага активности
type SetIsActiveRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	IsActive *bool  `json:"is_active" binding:"required"`
}

// SetIsActiveResponse - ответ при успешном обновлении
type SetIsActiveResponse struct {
	User models.User `json:"user"`
}

// GetReviewResponse - ответ при получении PR'ов пользователя
type GetReviewResponse struct {
	UserID       string                    `json:"user_id"`
	PullRequests []models.PullRequestShort `json:"pull_requests"`
}

// ErrorResponse соответствует схеме ErrorResponse из OpenAPI
type ErrorResponse = models.ErrorResponse
