package pullrequests

import (
	"net/http"

	"github.com/aabbuukkaarr8/PRService/internal/api/models"
	"github.com/gin-gonic/gin"
)

// ReassignReviewer обрабатывает POST /pullRequest/reassign
func (h *Handler) ReassignReviewer(c *gin.Context) {
	var req ReassignReviewerRequest

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
	resultPR, replacedBy, err := h.service.ReassignReviewer(c.Request.Context(), req.PullRequestID, req.OldUserID)
	if err != nil {
		// Проверяем тип ошибки
		if err.Error() == "NOT_FOUND" || err.Error() == "pullrequests: NOT_FOUND" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: struct {
					Code    models.ErrorResponseErrorCode `json:"code"`
					Message string                        `json:"message"`
				}{
					Code:    models.NOTFOUND,
					Message: "PR or user not found",
				},
			})
			return
		}

		if err.Error() == "PR_MERGED" || err.Error() == "pullrequests: PR_MERGED" {
			c.JSON(http.StatusConflict, models.ErrorResponse{
				Error: struct {
					Code    models.ErrorResponseErrorCode `json:"code"`
					Message string                        `json:"message"`
				}{
					Code:    models.PRMERGED,
					Message: "cannot reassign on merged PR",
				},
			})
			return
		}

		if err.Error() == "NOT_ASSIGNED" || err.Error() == "pullrequests: NOT_ASSIGNED" {
			c.JSON(http.StatusConflict, models.ErrorResponse{
				Error: struct {
					Code    models.ErrorResponseErrorCode `json:"code"`
					Message string                        `json:"message"`
				}{
					Code:    models.NOTASSIGNED,
					Message: "reviewer is not assigned to this PR",
				},
			})
			return
		}

		if err.Error() == "NO_CANDIDATE" || err.Error() == "pullrequests: NO_CANDIDATE" {
			c.JSON(http.StatusConflict, models.ErrorResponse{
				Error: struct {
					Code    models.ErrorResponseErrorCode `json:"code"`
					Message string                        `json:"message"`
				}{
					Code:    models.NOCANDIDATE,
					Message: "no active replacement candidate in team",
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
	c.JSON(http.StatusOK, ReassignReviewerResponse{
		PR:         handlerPR,
		ReplacedBy: replacedBy,
	})
}
