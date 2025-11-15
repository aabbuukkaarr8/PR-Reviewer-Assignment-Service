package pullrequests

import (
	"net/http"

	"github.com/aabbuukkaarr8/PRService/internal/service/pullrequests"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreatePullRequest(c *gin.Context) {
	var req CreatePullRequestRequest

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

	resultPR, err := h.service.CreatePullRequest(c.Request.Context(), req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		// Проверяем тип ошибки
		if err.Error() == "PR_EXISTS" || err.Error() == "pullrequests: PR_EXISTS" {
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: ErrorDetail{
					Code:    "PR_EXISTS",
					Message: "PR id already exists",
				},
			})
			return
		}

		if err.Error() == "NOT_FOUND" || err.Error() == "pullrequests: NOT_FOUND" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: ErrorDetail{
					Code:    "NOT_FOUND",
					Message: "author or team not found",
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
	c.JSON(http.StatusCreated, CreatePullRequestResponse{
		PR: handlerPR,
	})
}

// toHandlerPullRequest конвертирует service.PullRequest в handler.PullRequest
func toHandlerPullRequest(s pullrequests.PullRequest) PullRequest {
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
