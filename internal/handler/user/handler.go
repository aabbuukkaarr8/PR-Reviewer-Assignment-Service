package user

type Handler struct {
	service ServiceUser
}

func NewHandler(service ServiceUser) *Handler {
	return &Handler{
		service: service,
	}
}
