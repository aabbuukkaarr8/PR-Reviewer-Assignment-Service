package pullrequest

import (
	"errors"
	"net/http"

	"github.com/aabbuukkaarr8/PRService/internal/api"
	"github.com/aabbuukkaarr8/PRService/internal/api/models"
	prsrv "github.com/aabbuukkaarr8/PRService/internal/service/pullrequest"
	"github.com/gin-gonic/gin"
)

func (h *Handler) MergePullRequest(c *gin.Context) {
	var req MergeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		api.SendError(c, http.StatusBadRequest, api.Error{
			Code:    "INVALID_REQUEST",
			Message: err.Error(),
		})
		return
	}

	resultPR, err := h.service.MergePullRequest(c.Request.Context(), req.PullRequestID)
	if err != nil {
		switch {
		case errors.Is(err, prsrv.ErrNotFound):
			api.SendError(c, http.StatusNotFound, api.Error{
				Code:    models.NOTFOUND,
				Message: "PR not found",
			})
		default:
			api.SendError(c, http.StatusInternalServerError, api.Error{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			})
		}
		return
	}

	handlerPR := toHandlerPullRequest(resultPR)

	api.SendOk(c, MergePullRequestResponse{
		PR: handlerPR,
	})
}
