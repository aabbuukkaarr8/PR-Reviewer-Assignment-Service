package users

import (
	"net/http"

	"github.com/aabbuukkaarr8/PRService/internal/service/users"
	"github.com/gin-gonic/gin"
)

// SetIsActive обрабатывает POST /users/setIsActive
func (h *Handler) SetIsActive(c *gin.Context) {
	var req SetIsActiveRequest

	// Парсим JSON в структуру handler DTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
		return
	}

	// Передаем напрямую в service
	// req.IsActive - указатель, разыменовываем (binding:"required" гарантирует что не nil)
	resultUser, err := h.service.SetIsActive(c.Request.Context(), req.UserID, *req.IsActive)
	if err != nil {
		// Проверяем тип ошибки
		if err.Error() == "NOT_FOUND" || err.Error() == "user: NOT_FOUND" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: ErrorDetail{
					Code:    "NOT_FOUND",
					Message: "user not found",
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

	// Конвертируем service.User в handler.User
	handlerUser := toHandlerUser(resultUser)

	// Успешный ответ
	c.JSON(http.StatusOK, SetIsActiveResponse{
		User: handlerUser,
	})
}

// toHandlerUser конвертирует service.User в handler.User
func toHandlerUser(s users.User) User {
	return User{
		UserID:   s.UserID,
		Username: s.Username,
		TeamName: s.TeamName,
		IsActive: s.IsActive,
	}
}
