package user

import (
	"errors"
	"net/http"

	"github.com/aabbuukkaarr8/PRService/internal/api"
	"github.com/aabbuukkaarr8/PRService/internal/api/models"
	usersrv "github.com/aabbuukkaarr8/PRService/internal/service/user"
	"github.com/gin-gonic/gin"
)

func (h *Handler) SetIsActive(c *gin.Context) {
	var req SetIsActiveRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		api.SendError(c, http.StatusBadRequest, api.Error{
			Code:    "INVALID_REQUEST",
			Message: err.Error(),
		})
		return
	}

	resultUser, err := h.service.SetIsActive(c.Request.Context(), req.UserID, *req.IsActive)
	if err != nil {
		switch {
		case errors.Is(err, usersrv.ErrUserNotFound):
			api.SendError(c, http.StatusNotFound, api.Error{
				Code:    models.NOTFOUND,
				Message: "user not found",
			})
		default:
			api.SendError(c, http.StatusInternalServerError, api.Error{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			})
		}
		return
	}

	var handlerUser User
	handlerUser.FillFromService(resultUser)

	api.SendOk(c, SetIsActiveResponse{
		User: handlerUser,
	})
}
