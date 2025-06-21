package bootstrap

import (
	"github.com/craftaholic/insider/internal/controller"
	// "github.com/craftaholic/insider/internal/shared/env"
	// "github.com/craftaholic/insider/internal/shared/log"
)

type Application struct {
	// Controller/Handler Layer
	HealthController            *controller.HealthController
}

func App() Application {
	// logger := log.BaseLogger.WithFields("bootstrap", "App")

	// Initiate the application
	app := &Application{}

	return *app
}

func (app *Application) CloseDBConnection() {
	// Close the connection
}
