package team

import (
	"errors"
	"net/http"

	"github.com/aabbuukkaarr8/PRService/internal/api"
	"github.com/aabbuukkaarr8/PRService/internal/api/models"
	teamsrv "github.com/aabbuukkaarr8/PRService/internal/service/team"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetTeam(c *gin.Context) {
	teamName := c.Query("team_name")
	if teamName == "" {
		api.SendError(c, http.StatusBadRequest, api.Error{
			Code:    "INVALID_REQUEST",
			Message: "team_name parameter is required",
		})
		return
	}

	resultTeam, err := h.service.GetTeam(c.Request.Context(), teamName)
	if err != nil {
		switch {
		case errors.Is(err, teamsrv.ErrTeamNotFound):
			api.SendError(c, http.StatusNotFound, api.Error{
				Code:    models.NOTFOUND,
				Message: "team not found",
			})
		default:
			h.logger.WithError(err).WithField("team_name", teamName).Error("Failed to get team")
			api.SendError(c, http.StatusInternalServerError, api.Error{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			})
		}
		return
	}

	var handlerTeam Team
	handlerTeam.FillFromService(resultTeam)

	api.SendOk(c, GetTeamResponse{
		Team: handlerTeam,
	})
}
