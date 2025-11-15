package main

import (
	"flag"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/aabbuukkaarr8/PRService/internal/apiserver"
	pullrequests3 "github.com/aabbuukkaarr8/PRService/internal/handler/pullrequests"
	team3 "github.com/aabbuukkaarr8/PRService/internal/handler/team"
	users3 "github.com/aabbuukkaarr8/PRService/internal/handler/users"
	"github.com/aabbuukkaarr8/PRService/internal/repository/pullrequests"
	"github.com/aabbuukkaarr8/PRService/internal/repository/team"
	"github.com/aabbuukkaarr8/PRService/internal/repository/users"
	pullrequests2 "github.com/aabbuukkaarr8/PRService/internal/service/pullrequests"
	team2 "github.com/aabbuukkaarr8/PRService/internal/service/team"
	users2 "github.com/aabbuukkaarr8/PRService/internal/service/users"
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
	teamRepo := team.NewRepository(db)
	userRepo := users.NewRepository(db)
	prRepo := pullrequests.NewRepository(db)

	//srv
	teamSrv := team2.NewService(teamRepo)
	userSrv := users2.NewService(userRepo)
	prSrv := pullrequests2.NewService(prRepo)

	//handler
	teamHandler := team3.NewHandler(teamSrv)
	userHandler := users3.NewHandler(userSrv)
	prHandler := pullrequests3.NewHandler(prSrv)

	s := apiserver.New(config)
	s.ConfigureRouter(teamHandler, userHandler, prHandler)

	if err := s.Run(); err != nil {
		panic(err)
	}
}
