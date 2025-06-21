package main

import (
	"net/http"

	"github.com/craftaholic/insider/internal/api/route"
	"github.com/craftaholic/insider/internal/bootstrap"
	"github.com/craftaholic/insider/internal/shared/env"
	"github.com/craftaholic/insider/internal/shared/log"
)

func init() {
	log.Init()
	env.LoadEnv()

	log.BaseLogger.Info("Starting the application - Author: Tommy Tran - tommytrandt.work@gmail.com")
}

func main() {
	app := bootstrap.App()

	defer app.CloseDBConnection()

	r := route.SetupRoute(app)

	log.BaseLogger.Info("Starting server...", "on port", env.Env.ServerAddress)

	err := http.ListenAndServe(":"+env.Env.ServerAddress, r)
	if err != nil {
		log.BaseLogger.Error("Server bootstraping error", "error", err)
	}

	log.BaseLogger.Info("Stopping server...", "on port", env.Env.ServerAddress)
}
