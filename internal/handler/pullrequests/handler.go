package pullrequests

// Handler структура для обработчиков pull requests
type Handler struct {
	service PullRequestService
}

// NewHandler создает новый Handler
func NewHandler(service PullRequestService) *Handler {
	return &Handler{
		service: service,
	}
}
