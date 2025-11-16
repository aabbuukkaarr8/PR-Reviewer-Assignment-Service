package pullrequest

import (
	"errors"
	"net/http"

	"github.com/aabbuukkaarr8/PRService/internal/api"
	"github.com/aabbuukkaarr8/PRService/internal/api/models"
	prsrv "github.com/aabbuukkaarr8/PRService/internal/service/pullrequest"
	"github.com/gin-gonic/gin"
)

type CreateRequest struct {
	AuthorID        string `json:"author_id" binding:"required"`
	PullRequestID   string `json:"pull_request_id" binding:"required"`
	PullRequestName string `json:"pull_request_name" binding:"required"`
}

func (h *Handler) CreatePullRequest(c *gin.Context) {
	var req CreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		api.SendError(c, http.StatusBadRequest, api.Error{
			Code:    "INVALID_REQUEST",
			Message: err.Error(),
		})
		return
	}

	reqToSrv := prsrv.CreatePullRequest{
		AuthorId:        req.AuthorID,
		PullRequestId:   req.PullRequestID,
		PullRequestName: req.PullRequestName,
	}

	resultPR, err := h.service.CreatePullRequest(c.Request.Context(), reqToSrv)
	if err != nil {
		switch {
		case errors.Is(err, prsrv.ErrPRExists):
			api.SendError(c, http.StatusConflict, api.Error{
				Code:    models.PREXISTS,
				Message: "PR id already exists",
			})
		case errors.Is(err, prsrv.ErrNotFound):
			api.SendError(c, http.StatusNotFound, api.Error{
				Code:    models.NOTFOUND,
				Message: "author or team not found",
			})
		default:
			h.logger.WithError(err).WithField("pull_request_id", req.PullRequestID).Error("Failed to create PR")
			api.SendError(c, http.StatusInternalServerError, api.Error{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			})
		}
		return
	}

	handlerPR := toHandlerPullRequest(resultPR)

	api.SendCreated(c, CreatePullRequestResponse{
		PR: handlerPR,
	})
}
