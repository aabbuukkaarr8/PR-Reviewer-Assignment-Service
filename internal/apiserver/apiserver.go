package apiserver

import (
	"net/http"

	"github.com/aabbuukkaarr8/PRService/internal/handler/pullrequest"
	"github.com/aabbuukkaarr8/PRService/internal/handler/user"

	"github.com/aabbuukkaarr8/PRService/internal/handler/team"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *gin.Engine
}

func New(config *Config) *APIServer {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	return &APIServer{
		config: config,
		logger: logger,
		router: gin.Default(),
	}

}

func (s *APIServer) Run() error {
	if err := s.configLogger(); err != nil {
		return err
	}

	s.logger.Info("Starting API server")
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) configLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func (s *APIServer) ConfigureRouter(teamHandler *team.Handler, usersHandler *user.Handler, prHandler *pullrequest.Handler) {
	s.router.POST("/team/add", teamHandler.CreateTeam)
	s.router.GET("/team/get", teamHandler.GetTeam)
	s.router.POST("/users/setIsActive", usersHandler.SetIsActive)
	s.router.GET("/users/getReview", usersHandler.GetReview)
	s.router.POST("/pullRequest/create", prHandler.CreatePullRequest)
	s.router.POST("/pullRequest/merge", prHandler.MergePullRequest)
	s.router.POST("/pullRequest/reassign", prHandler.ReassignReviewer)
}

func (s *APIServer) GetRouter() *gin.Engine {
	return s.router
}
