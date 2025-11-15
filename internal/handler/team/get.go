package team

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetTeam обрабатывает GET /team/get
func (h *Handler) GetTeam(c *gin.Context) {
	// Получаем team_name из query параметров
	teamName := c.Query("team_name")
	if teamName == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
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
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: ErrorDetail{
					Code:    "NOT_FOUND",
					Message: "team not found",
				},
			})
			return
		}

		// Другие ошибки
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	// Конвертируем service.Team в handler.Team
	handlerTeam := toHandlerTeam(resultTeam)

	// Успешный ответ
	c.JSON(http.StatusOK, handlerTeam)
}
