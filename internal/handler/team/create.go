package team

import (
	"net/http"

	"github.com/aabbuukkaarr8/PRService/internal/api/models"
	"github.com/aabbuukkaarr8/PRService/internal/service/team"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateTeam(c *gin.Context) {
	var req CreateTeamRequest

	// Парсим JSON в структуру handler DTO
	if err := c.ShouldBindJSON(&req.Team); err != nil {
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

	// Конвертируем models.Team в service.Team
	serviceTeam := toServiceTeam(req.Team)

	// Передаем в service
	resultTeam, err := h.service.CreateTeam(c.Request.Context(), serviceTeam)
	if err != nil {
		// Проверяем тип ошибки
		if err.Error() == "TEAM_EXISTS" || err.Error() == "team: TEAM_EXISTS" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: struct {
					Code    models.ErrorResponseErrorCode `json:"code"`
					Message string                        `json:"message"`
				}{
					Code:    models.TEAMEXISTS,
					Message: "team_name already exists",
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

	// Конвертируем service.Team обратно в models.Team
	handlerTeam := toHandlerTeam(resultTeam)

	c.JSON(http.StatusCreated, CreateTeamResponse{
		Team: handlerTeam,
	})
}

// toServiceTeam конвертирует models.Team в service.Team
func toServiceTeam(h models.Team) team.Team {
	members := make([]team.TeamMember, len(h.Members))
	for i, m := range h.Members {
		members[i] = team.TeamMember{
			UserID:   m.UserId,
			Username: m.Username,
			IsActive: m.IsActive,
		}
	}

	return team.Team{
		TeamName: h.TeamName,
		Members:  members,
	}
}

// toHandlerTeam конвертирует service.Team в models.Team
func toHandlerTeam(s team.Team) models.Team {
	members := make([]models.TeamMember, len(s.Members))
	for i, m := range s.Members {
		members[i] = models.TeamMember{
			UserId:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		}
	}

	return models.Team{
		TeamName: s.TeamName,
		Members:  members,
	}
}
