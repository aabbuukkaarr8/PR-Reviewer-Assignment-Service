package team

// TeamMember соответствует схеме TeamMember из OpenAPI
type TeamMember struct {
	UserID   string `json:"user_id" binding:"required"`
	Username string `json:"username" binding:"required"`
	IsActive bool   `json:"is_active" binding:"required"`
}

// Team соответствует схеме Team из OpenAPI
type Team struct {
	TeamName string       `json:"team_name" binding:"required"`
	Members  []TeamMember `json:"members" binding:"required"`
}

// CreateTeamRequest - запрос на создание команды
type CreateTeamRequest struct {
	Team
}

// CreateTeamResponse - ответ при успешном создании команды
type CreateTeamResponse struct {
	Team Team `json:"team"`
}

// ErrorResponse соответствует схеме ErrorResponse из OpenAPI
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail детали ошибки
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
