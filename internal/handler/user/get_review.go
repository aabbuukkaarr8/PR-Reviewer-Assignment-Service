package user

import (
	"errors"
	"net/http"

	"github.com/aabbuukkaarr8/PRService/internal/api"
	"github.com/aabbuukkaarr8/PRService/internal/api/models"
	usersrv "github.com/aabbuukkaarr8/PRService/internal/service/user"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetReview(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		api.SendError(c, http.StatusBadRequest, api.Error{
			Code:    "INVALID_REQUEST",
			Message: "user_id query parameter is required",
		})
		return
	}

	resultPRs, err := h.service.GetReview(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, usersrv.ErrUserNotFound):
			api.SendError(c, http.StatusNotFound, api.Error{
				Code:    models.NOTFOUND,
				Message: "user not found",
			})
		default:
			api.SendError(c, http.StatusInternalServerError, api.Error{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			})
		}
		return
	}

	handlerPRs := make([]PullRequestShort, len(resultPRs))
	for i, pr := range resultPRs {
		handlerPRs[i].FillFromService(pr)
	}

	api.SendOk(c, GetReviewResponse{
		UserID:       userID,
		PullRequests: handlerPRs,
	})
}
