package pullrequest

import (
	"errors"
	"net/http"

	"github.com/aabbuukkaarr8/PRService/internal/api"
	"github.com/aabbuukkaarr8/PRService/internal/api/models"
	prsrv "github.com/aabbuukkaarr8/PRService/internal/service/pullrequest"
	"github.com/gin-gonic/gin"
)

func (h *Handler) ReassignReviewer(c *gin.Context) {
	var req ReassignReviewerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		api.SendError(c, http.StatusBadRequest, api.Error{
			Code:    "INVALID_REQUEST",
			Message: err.Error(),
		})
		return
	}

	resultPR, replacedBy, err := h.service.ReassignReviewer(c.Request.Context(), req.PullRequestID, req.OldUserID)
	if err != nil {
		switch {
		case errors.Is(err, prsrv.ErrNotFound):
			api.SendError(c, http.StatusNotFound, api.Error{
				Code:    models.NOTFOUND,
				Message: "PR or user not found",
			})
		case errors.Is(err, prsrv.ErrPRMerged):
			api.SendError(c, http.StatusConflict, api.Error{
				Code:    models.PRMERGED,
				Message: "cannot reassign on merged PR",
			})
		case errors.Is(err, prsrv.ErrNotAssigned):
			api.SendError(c, http.StatusConflict, api.Error{
				Code:    models.NOTASSIGNED,
				Message: "reviewer is not assigned to this PR",
			})
		case errors.Is(err, prsrv.ErrNoCandidate):
			api.SendError(c, http.StatusConflict, api.Error{
				Code:    models.NOCANDIDATE,
				Message: "no active replacement candidate in team",
			})
		default:
			h.logger.WithError(err).WithFields(map[string]interface{}{
				"pull_request_id": req.PullRequestID,
				"old_reviewer_id": req.OldUserID,
			}).Error("Failed to reassign reviewer")
			api.SendError(c, http.StatusInternalServerError, api.Error{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			})
		}
		return
	}

	handlerPR := toHandlerPullRequest(resultPR)

	api.SendOk(c, ReassignReviewerResponse{
		PR:         handlerPR,
		ReplacedBy: replacedBy,
	})
}
