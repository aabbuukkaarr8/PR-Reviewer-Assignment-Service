package pullrequests

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MergePullRequest обрабатывает POST /pullRequest/merge
func (h *Handler) MergePullRequest(c *gin.Context) {
	var req MergePullRequestRequest

	// Парсим JSON в структуру handler DTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
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
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: ErrorDetail{
					Code:    "NOT_FOUND",
					Message: "PR not found",
				},
			})
			return
		}

		// Другие ошибки
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	// Конвертируем service.PullRequest в handler.PullRequest
	handlerPR := toHandlerPullRequest(resultPR)

	// Успешный ответ
	c.JSON(http.StatusOK, MergePullRequestResponse{
		PR: handlerPR,
	})
}
