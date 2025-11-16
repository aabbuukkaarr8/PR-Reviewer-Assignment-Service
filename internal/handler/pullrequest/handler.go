package pullrequest

type Handler struct {
	service ServicePR
}

func NewHandler(service ServicePR) *Handler {
	return &Handler{
		service: service,
	}
}
