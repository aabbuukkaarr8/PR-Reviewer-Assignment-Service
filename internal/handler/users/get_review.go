package users

import (
	"net/http"

	"github.com/aabbuukkaarr8/PRService/internal/service/users"
	"github.com/gin-gonic/gin"
)

// GetReview обрабатывает GET /users/getReview
func (h *Handler) GetReview(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "INVALID_REQUEST",
				Message: "user_id query parameter is required",
			},
		})
		return
	}

	// Получаем PR'ы пользователя из service
	resultPRs, err := h.service.GetReview(c.Request.Context(), userID)
	if err != nil {
		// Проверяем тип ошибки
		if err.Error() == "NOT_FOUND" || err.Error() == "users: NOT_FOUND" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: ErrorDetail{
					Code:    "NOT_FOUND",
					Message: "user not found",
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

	// Конвертируем service.PullRequestShort в handler.PullRequestShort
	handlerPRs := make([]PullRequestShort, len(resultPRs))
	for i, pr := range resultPRs {
		handlerPRs[i] = toHandlerPullRequestShort(pr)
	}

	// Успешный ответ
	c.JSON(http.StatusOK, GetReviewResponse{
		UserID:       userID,
		PullRequests: handlerPRs,
	})
}

// toHandlerPullRequestShort конвертирует service.PullRequestShort в handler.PullRequestShort
func toHandlerPullRequestShort(s users.PullRequestShort) PullRequestShort {
	return PullRequestShort{
		PullRequestID:   s.PullRequestID,
		PullRequestName: s.PullRequestName,
		AuthorID:        s.AuthorID,
		Status:          s.Status,
	}
}
