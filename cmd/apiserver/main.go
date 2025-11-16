package main

import (
	"flag"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/aabbuukkaarr8/PRService/internal/apiserver"
	prapi "github.com/aabbuukkaarr8/PRService/internal/handler/pullrequest"
	teamapi "github.com/aabbuukkaarr8/PRService/internal/handler/team"
	userapi "github.com/aabbuukkaarr8/PRService/internal/handler/user"
	prrepo "github.com/aabbuukkaarr8/PRService/internal/repository/pullrequest"
	teamrepo "github.com/aabbuukkaarr8/PRService/internal/repository/team"
	userrepo "github.com/aabbuukkaarr8/PRService/internal/repository/user"
	prsrv "github.com/aabbuukkaarr8/PRService/internal/service/pullrequest"
	teamsrv "github.com/aabbuukkaarr8/PRService/internal/service/team"
	usersrv "github.com/aabbuukkaarr8/PRService/internal/service/user"
	"github.com/aabbuukkaarr8/PRService/internal/store"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func main() {
	flag.Parse()
	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	// Переменная окружения DATABASE_URL имеет приоритет над конфигом
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		config.Store.DatabaseURL = dbURL
	}

	db := store.New()
	err = db.Open(config.Store.DatabaseURL)
	if err != nil {
		log.Fatal(err)
		return
	}
	//repo
	teamRepo := teamrepo.NewRepository(db)
	userRepo := userrepo.NewRepository(db)
	prRepo := prrepo.NewRepository(db)

	//srv
	teamSrv := teamsrv.NewService(teamRepo)
	userSrv := usersrv.NewService(userRepo)
	prSrv := prsrv.NewService(prRepo)

	//handler
	teamHandler := teamapi.NewHandler(teamSrv)
	userHandler := userapi.NewHandler(userSrv)
	prHandler := prapi.NewHandler(prSrv)

	s := apiserver.New(config)
	s.ConfigureRouter(teamHandler, userHandler, prHandler)

	if err := s.Run(); err != nil {
		panic(err)
	}
}
