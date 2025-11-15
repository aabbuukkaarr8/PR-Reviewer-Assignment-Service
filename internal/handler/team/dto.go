package team

import (
	"github.com/aabbuukkaarr8/PRService/internal/api/models"
)

// CreateTeamRequest - запрос на создание команды
type CreateTeamRequest struct {
	models.Team
}

// CreateTeamResponse - ответ при успешном создании команды
type CreateTeamResponse struct {
	Team models.Team `json:"team"`
}

// GetTeamResponse - ответ при получении команды
type GetTeamResponse struct {
	Team models.Team `json:"team"`
}

// ErrorResponse соответствует схеме ErrorResponse из OpenAPI
type ErrorResponse = models.ErrorResponse
