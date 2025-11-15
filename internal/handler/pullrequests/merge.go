package pullrequests

import (
	"net/http"

	"github.com/aabbuukkaarr8/PRService/internal/api/models"
	"github.com/aabbuukkaarr8/PRService/internal/service/pullrequests"
	"github.com/gin-gonic/gin"
)

// toHandlerPullRequest конвертирует service.PullRequest в models.PullRequest
func toHandlerPullRequest(s pullrequests.PullRequest) models.PullRequest {
	return models.PullRequest{
		PullRequestId:     s.PullRequestID,
		PullRequestName:   s.PullRequestName,
		AuthorId:          s.AuthorID,
		Status:            models.PullRequestStatus(s.Status),
		AssignedReviewers: s.AssignedReviewers,
		CreatedAt:         s.CreatedAt,
		MergedAt:          s.MergedAt,
	}
}

// MergePullRequest обрабатывает POST /pullRequest/merge
func (h *Handler) MergePullRequest(c *gin.Context) {
	var req MergePullRequestRequest

	// Парсим JSON в структуру handler DTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: struct {
				Code    models.ErrorResponseErrorCode `json:"code"`
				Message string                        `json:"message"`
			}{
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
		return
	}

	// Передаем напрямую в service
	resultPR, err := h.service.MergePullRequest(c.Request.Context(), req.PullRequestID)
	if err != nil {
		// Проверяем тип ошибки
		if err.Error() == "NOT_FOUND" || err.Error() == "pullrequests: NOT_FOUND" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: struct {
					Code    models.ErrorResponseErrorCode `json:"code"`
					Message string                        `json:"message"`
				}{
					Code:    models.NOTFOUND,
					Message: "PR not found",
				},
			})
			return
		}

		// Другие ошибки
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: struct {
				Code    models.ErrorResponseErrorCode `json:"code"`
				Message string                        `json:"message"`
			}{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	// Конвертируем service.PullRequest в models.PullRequest
	handlerPR := toHandlerPullRequest(resultPR)

	// Успешный ответ
	c.JSON(http.StatusOK, MergePullRequestResponse{
		PR: handlerPR,
	})
}
