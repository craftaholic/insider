package bootstrap

import (
	"context"
	"fmt"

	"github.com/craftaholic/insider/internal/controller"
	"github.com/craftaholic/insider/internal/repository"
	"github.com/craftaholic/insider/internal/usecase"
	"github.com/go-redis/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"

	"github.com/craftaholic/insider/internal/shared/config"
	"github.com/craftaholic/insider/internal/shared/log"
)

type Application struct {
	// Controller/Handler Layer
	HealthController  *controller.HealthController
	MessageController *controller.MessageController
}

func App() Application {
	logger := log.BaseLogger.WithFields("bootstrap", "App")

	// Initiate the application
	app := &Application{}

	// Init Infra Layer
	// Init DB conection
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Env.DBHost,
		config.Env.DBUser,
		config.Env.DBPassword,
		config.Env.DBName,
		config.Env.DBPort,
		config.Env.DBSslMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlog.Default.LogMode(gormlog.Error),
	})
	if err != nil {
		logger.Fatal("Failed to connect to database:", err)
	}

	// Init Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%s",
			config.Env.RedisHost,
			config.Env.RedisPort,
		),
	})

	// Init Repository Layer
	messageRepository := repository.NewMessageRepository(db)
	cacheRepository := repository.NewCacheRepository(redisClient)
	notificationService := repository.NewNotificationService(config.Env.WebhookAuthKey, config.Env.WebhookURL)

	// Init Usecase Layer
	messageUsecase := usecase.NewMessageUsecase(messageRepository, cacheRepository, notificationService)

	// Init Controller
	app.HealthController = controller.NewHealthController()
	app.MessageController = controller.NewMessageController(messageUsecase)

	// Start the automation for sending message

	// Execute the start automated sending in background context
	err = messageUsecase.StartAutomatedSending(context.Background())
	if err != nil {
		logger.Fatal("Error starting automated sending function")
	}

	return *app
}

func (app *Application) CloseDBConnection() {
	// Close the connection
}
