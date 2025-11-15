package pullrequests

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ReassignReviewer обрабатывает POST /pullRequest/reassign
func (h *Handler) ReassignReviewer(c *gin.Context) {
	var req ReassignReviewerRequest

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
	resultPR, replacedBy, err := h.service.ReassignReviewer(c.Request.Context(), req.PullRequestID, req.OldUserID)
	if err != nil {
		// Проверяем тип ошибки
		if err.Error() == "NOT_FOUND" || err.Error() == "pullrequests: NOT_FOUND" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: ErrorDetail{
					Code:    "NOT_FOUND",
					Message: "PR or user not found",
				},
			})
			return
		}

		if err.Error() == "PR_MERGED" || err.Error() == "pullrequests: PR_MERGED" {
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: ErrorDetail{
					Code:    "PR_MERGED",
					Message: "cannot reassign on merged PR",
				},
			})
			return
		}

		if err.Error() == "NOT_ASSIGNED" || err.Error() == "pullrequests: NOT_ASSIGNED" {
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: ErrorDetail{
					Code:    "NOT_ASSIGNED",
					Message: "reviewer is not assigned to this PR",
				},
			})
			return
		}

		if err.Error() == "NO_CANDIDATE" || err.Error() == "pullrequests: NO_CANDIDATE" {
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: ErrorDetail{
					Code:    "NO_CANDIDATE",
					Message: "no active replacement candidate in team",
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
	c.JSON(http.StatusOK, ReassignReviewerResponse{
		PR:         handlerPR,
		ReplacedBy: replacedBy,
	})
}
