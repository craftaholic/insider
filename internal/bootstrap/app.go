package bootstrap

import (
	"github.com/craftaholic/insider/internal/controller"
	"github.com/craftaholic/insider/internal/repository"
	"github.com/craftaholic/insider/internal/usecase"
	"github.com/go-redis/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"

	// "github.com/craftaholic/insider/internal/shared/env"
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
	dsn := "host=localhost user=postgres password=postgres123 dbname=message_system port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlog.Default.LogMode(gormlog.Info),
	})
	if err != nil {
		logger.Fatal("Failed to connect to database:", err)
	}

	// Init Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Init Repository Layer
	messageRepository := repository.NewMessageRepository(db)
	cacheRepository := repository.NewCacheRepository(redisClient)
	notificationService := repository.NewNotificationService("", "")

	// Init Usecase Layer
	messageUsecase := usecase.NewMessageUsecase(messageRepository, cacheRepository, notificationService)

	// Init Controller
	app.HealthController = controller.NewHealthController()
	app.MessageController = controller.NewMessageController(messageUsecase)

	return *app
}

func (app *Application) CloseDBConnection() {
	// Close the connection
}
