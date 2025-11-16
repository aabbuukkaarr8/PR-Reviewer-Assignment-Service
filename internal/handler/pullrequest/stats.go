package pullrequest

import (
	"net/http"

	"github.com/aabbuukkaarr8/PRService/internal/api"
	"github.com/gin-gonic/gin"
)

type ReviewerStats struct {
	UserID           string `json:"user_id"`
	Username         string `json:"username"`
	TeamName         string `json:"team_name"`
	AssignmentsCount int    `json:"assignments_count"`
}

type PRStats struct {
	TotalPRs  int `json:"total_prs"`
	OpenPRs   int `json:"open_prs"`
	MergedPRs int `json:"merged_prs"`
}

type StatsResponse struct {
	PRStats       PRStats         `json:"pr_stats"`
	ReviewerStats []ReviewerStats `json:"reviewer_stats"`
}

func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.service.GetStats(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Error("Failed to get statistics")
		api.SendError(c, http.StatusInternalServerError, api.Error{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	handlerStats := StatsResponse{
		PRStats: PRStats{
			TotalPRs:  stats.PRStats.TotalPRs,
			OpenPRs:   stats.PRStats.OpenPRs,
			MergedPRs: stats.PRStats.MergedPRs,
		},
		ReviewerStats: make([]ReviewerStats, len(stats.ReviewerStats)),
	}

	for i, rs := range stats.ReviewerStats {
		handlerStats.ReviewerStats[i] = ReviewerStats{
			UserID:           rs.UserID,
			Username:         rs.Username,
			TeamName:         rs.TeamName,
			AssignmentsCount: rs.AssignmentsCount,
		}
	}

	api.SendOk(c, handlerStats)
}
