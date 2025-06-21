package main

import (
	"net/http"
	"time"

	"github.com/craftaholic/insider/internal/api/route"
	"github.com/craftaholic/insider/internal/bootstrap"
	"github.com/craftaholic/insider/internal/shared/config"
	"github.com/craftaholic/insider/internal/shared/constant"
	"github.com/craftaholic/insider/internal/shared/log"
)

func main() {
	log.Init()
	config.LoadEnv()

	baseLogger := log.BaseLogger
	baseLogger.Info("Starting the application - Author: Tommy Tran - tommytrandt.work@gmail.com")
	app := bootstrap.App()

	defer app.CloseDBConnection()

	r := route.SetupRoute(app)

	baseLogger.Info("Starting server...", "on port", config.Env.ServerAddress)

	srv := &http.Server{
		Addr:         ":" + config.Env.ServerAddress,
		Handler:      r,
		ReadTimeout:  constant.DefaultTimeout * time.Second,
		WriteTimeout: constant.IdleTimeout * time.Second,
		IdleTimeout:  constant.DefaultTimeout * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		baseLogger.Error("Server bootstraping error", "error", err)
	}

	baseLogger.Info("Stopping server...", "on port", config.Env.ServerAddress)
}
