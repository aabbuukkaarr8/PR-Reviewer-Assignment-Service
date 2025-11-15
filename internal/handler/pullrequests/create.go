package pullrequests

import (
	"net/http"

	"github.com/aabbuukkaarr8/PRService/internal/api/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreatePullRequest(c *gin.Context) {
	var req CreatePullRequestRequest

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

	resultPR, err := h.service.CreatePullRequest(c.Request.Context(), req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		// Проверяем тип ошибки
		if err.Error() == "PR_EXISTS" || err.Error() == "pullrequests: PR_EXISTS" {
			c.JSON(http.StatusConflict, models.ErrorResponse{
				Error: struct {
					Code    models.ErrorResponseErrorCode `json:"code"`
					Message string                        `json:"message"`
				}{
					Code:    models.PREXISTS,
					Message: "PR id already exists",
				},
			})
			return
		}

		if err.Error() == "NOT_FOUND" || err.Error() == "pullrequests: NOT_FOUND" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: struct {
					Code    models.ErrorResponseErrorCode `json:"code"`
					Message string                        `json:"message"`
				}{
					Code:    models.NOTFOUND,
					Message: "author or team not found",
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
	c.JSON(http.StatusCreated, CreatePullRequestResponse{
		PR: handlerPR,
	})
}
