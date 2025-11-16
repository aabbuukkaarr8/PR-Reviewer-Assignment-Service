package api

import (
	"encoding/json"
	"net/http"

	"github.com/aabbuukkaarr8/PRService/internal/api/models"
	"github.com/gin-gonic/gin"
)

type Error struct {
	Code    models.ErrorResponseErrorCode `json:"code"`
	Message string                        `json:"message"`
}

func SendError(c *gin.Context, statusCode int, response Error) {
	c.JSON(statusCode, models.ErrorResponse{Error: response})
}

func SendOk(c *gin.Context, response any) {
	_, err := json.Marshal(response)
	if err != nil {
		SendError(c, http.StatusInternalServerError, Error{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func SendCreated(c *gin.Context, response any) {
	_, err := json.Marshal(response)
	if err != nil {
		SendError(c, http.StatusInternalServerError, Error{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}
