package user

import (
	"github.com/sirupsen/logrus"
)

type Handler struct {
	service ServiceUser
	logger  *logrus.Logger
}

func NewHandler(service ServiceUser, logger *logrus.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}
