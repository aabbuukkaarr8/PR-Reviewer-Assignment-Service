package team

import (
	"net/http"

	"github.com/aabbuukkaarr8/PRService/internal/api/models"
	"github.com/gin-gonic/gin"
)

// GetTeam обрабатывает GET /team/get
func (h *Handler) GetTeam(c *gin.Context) {
	// Получаем team_name из query параметров
	teamName := c.Query("team_name")
	if teamName == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: struct {
				Code    models.ErrorResponseErrorCode `json:"code"`
				Message string                        `json:"message"`
			}{
				Code:    "INVALID_REQUEST",
				Message: "team_name parameter is required",
			},
		})
		return
	}

	// Передаем в service
	resultTeam, err := h.service.GetTeam(c.Request.Context(), teamName)
	if err != nil {
		// Проверяем тип ошибки
		if err.Error() == "NOT_FOUND" || err.Error() == "team: NOT_FOUND" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: struct {
					Code    models.ErrorResponseErrorCode `json:"code"`
					Message string                        `json:"message"`
				}{
					Code:    models.NOTFOUND,
					Message: "team not found",
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

	// Конвертируем service.Team в models.Team
	handlerTeam := toHandlerTeam(resultTeam)

	// Успешный ответ
	c.JSON(http.StatusOK, GetTeamResponse{
		Team: handlerTeam,
	})
}
