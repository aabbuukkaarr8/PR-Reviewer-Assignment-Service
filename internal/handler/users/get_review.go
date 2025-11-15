package users

import (
	"net/http"

	"github.com/aabbuukkaarr8/PRService/internal/api/models"
	"github.com/aabbuukkaarr8/PRService/internal/service/users"
	"github.com/gin-gonic/gin"
)

// GetReview обрабатывает GET /users/getReview
func (h *Handler) GetReview(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: struct {
				Code    models.ErrorResponseErrorCode `json:"code"`
				Message string                        `json:"message"`
			}{
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
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: struct {
					Code    models.ErrorResponseErrorCode `json:"code"`
					Message string                        `json:"message"`
				}{
					Code:    models.NOTFOUND,
					Message: "user not found",
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

	// Конвертируем service.PullRequestShort в models.PullRequestShort
	handlerPRs := make([]models.PullRequestShort, len(resultPRs))
	for i, pr := range resultPRs {
		handlerPRs[i] = toHandlerPullRequestShort(pr)
	}

	// Успешный ответ
	c.JSON(http.StatusOK, GetReviewResponse{
		UserID:       userID,
		PullRequests: handlerPRs,
	})
}

// toHandlerPullRequestShort конвертирует service.PullRequestShort в models.PullRequestShort
func toHandlerPullRequestShort(s users.PullRequestShort) models.PullRequestShort {
	return models.PullRequestShort{
		PullRequestId:   s.PullRequestID,
		PullRequestName: s.PullRequestName,
		AuthorId:        s.AuthorID,
		Status:          models.PullRequestShortStatus(s.Status),
	}
}
