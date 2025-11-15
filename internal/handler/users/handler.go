package users

// Handler структура для обработчиков команд
type Handler struct {
	service UserService
}

// NewHandler создает новый Handler
func NewHandler(service UserService) *Handler {
	return &Handler{
		service: service,
	}
}
