package team

import (
	"errors"
	"net/http"

	"github.com/aabbuukkaarr8/PRService/internal/api"
	"github.com/aabbuukkaarr8/PRService/internal/api/models"
	teamsrv "github.com/aabbuukkaarr8/PRService/internal/service/team"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateTeam(c *gin.Context) {
	var req CreateTeamRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		api.SendError(c, http.StatusBadRequest, api.Error{
			Code:    "INVALID_REQUEST",
			Message: err.Error(),
		})
		return
	}

	serviceTeam := req.ToService()

	resultTeam, err := h.service.CreateTeam(c.Request.Context(), serviceTeam)
	if err != nil {
		switch {
		case errors.Is(err, teamsrv.ErrTeamExists):
			api.SendError(c, http.StatusBadRequest, api.Error{
				Code:    models.TEAMEXISTS,
				Message: "team_name already exists",
			})
		default:
			api.SendError(c, http.StatusInternalServerError, api.Error{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			})
		}
		return
	}

	var handlerTeam Team
	handlerTeam.FillFromService(resultTeam)

	api.SendCreated(c, CreateTeamResponse{
		Team: handlerTeam,
	})
}
