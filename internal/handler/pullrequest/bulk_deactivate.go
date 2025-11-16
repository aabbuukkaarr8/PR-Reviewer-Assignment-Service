package pullrequest

import (
	"net/http"

	"github.com/aabbuukkaarr8/PRService/internal/api"
	"github.com/gin-gonic/gin"
)

type BulkDeactivateRequest struct {
	TeamName string `json:"team_name" binding:"required"`
}

type ReassignedPR struct {
	PullRequestID string `json:"pull_request_id"`
	OldReviewerID string `json:"old_reviewer_id"`
	NewReviewerID string `json:"new_reviewer_id"`
}

type BulkDeactivateResponse struct {
	DeactivatedUserIDs []string       `json:"deactivated_user_ids"`
	ReassignedPRs      []ReassignedPR `json:"reassigned_prs"`
}

func (h *Handler) BulkDeactivateTeamUsers(c *gin.Context) {
	var req BulkDeactivateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		api.SendError(c, http.StatusBadRequest, api.Error{
			Code:    "INVALID_REQUEST",
			Message: err.Error(),
		})
		return
	}

	result, err := h.service.BulkDeactivateTeamUsers(c.Request.Context(), req.TeamName)
	if err != nil {
		h.logger.WithError(err).WithField("team_name", req.TeamName).Error("Failed to bulk deactivate team users")
		api.SendError(c, http.StatusInternalServerError, api.Error{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	handlerReassignedPRs := make([]ReassignedPR, len(result.ReassignedPRs))
	for i, pr := range result.ReassignedPRs {
		handlerReassignedPRs[i] = ReassignedPR{
			PullRequestID: pr.PullRequestID,
			OldReviewerID: pr.OldReviewerID,
			NewReviewerID: pr.NewReviewerID,
		}
	}

	api.SendOk(c, BulkDeactivateResponse{
		DeactivatedUserIDs: result.DeactivatedUserIDs,
		ReassignedPRs:      handlerReassignedPRs,
	})
}
