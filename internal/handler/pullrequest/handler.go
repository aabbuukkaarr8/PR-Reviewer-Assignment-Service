package pullrequest

import (
	"github.com/sirupsen/logrus"
)

type Handler struct {
	service ServicePR
	logger  *logrus.Logger
}

func NewHandler(service ServicePR, logger *logrus.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}
