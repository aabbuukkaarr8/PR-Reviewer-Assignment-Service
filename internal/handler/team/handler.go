package team

type Handler struct {
	service ServiceTeam
}

func NewHandler(service ServiceTeam) *Handler {
	return &Handler{
		service: service,
	}
}
