package team

// Handler структура для обработчиков команд
type Handler struct {
	service TeamService
}

// NewHandler создает новый Handler
func NewHandler(service TeamService) *Handler {
	return &Handler{
		service: service,
	}
}
