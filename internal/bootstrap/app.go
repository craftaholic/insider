package bootstrap

import (
	"github.com/craftaholic/insider/internal/controller"
	"github.com/craftaholic/insider/internal/domain"
	"github.com/craftaholic/insider/internal/repository"
	"github.com/craftaholic/insider/internal/usecase"
	// "github.com/craftaholic/insider/internal/shared/env"
	// "github.com/craftaholic/insider/internal/shared/log".
)

type Application struct {
	// Controller/Handler Layer
	HealthController  *controller.HealthController
	MessageController *controller.MessageController
}

func App() Application {
	// logger := log.BaseLogger.WithFields("bootstrap", "App")

	// Initiate the application
	app := &Application{}

	// Init Infra Layer
	
	
	// Init Repository Layer
	messageRepository := repository.NewMessageRepository()

	// Init Usecase Layer
	messageUsecase := usecase.NewMessageUsecase(messageRepository)

	// Init Controller
	app.HealthController = controller.NewHealthController()
	app.MessageController = controller.NewMessageController(messageUsecase)

	return *app
}

func (app *Application) CloseDBConnection() {
	// Close the connection
}
