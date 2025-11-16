package team

import (
	"github.com/sirupsen/logrus"
)

type Handler struct {
	service ServiceTeam
	logger  *logrus.Logger
}

func NewHandler(service ServiceTeam, logger *logrus.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}
