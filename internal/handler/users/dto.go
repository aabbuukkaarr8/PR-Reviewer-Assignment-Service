package users

// User соответствует схеме User из OpenAPI
type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

// SetIsActiveRequest - запрос на установку флага активности
type SetIsActiveRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	IsActive *bool  `json:"is_active" binding:"required"`
}

// SetIsActiveResponse - ответ при успешном обновлении
type SetIsActiveResponse struct {
	User User `json:"user"`
}

// PullRequestShort соответствует схеме PullRequestShort из OpenAPI
type PullRequestShort struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"` // OPEN or MERGED
}

// GetReviewResponse - ответ при получении PR'ов пользователя
type GetReviewResponse struct {
	UserID       string             `json:"user_id"`
	PullRequests []PullRequestShort `json:"pull_requests"`
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
