package bootstrap

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/craftaholic/insider/internal/controller"
	"github.com/craftaholic/insider/internal/repository"
	"github.com/craftaholic/insider/internal/usecase"
	"github.com/go-redis/redis"
	"github.com/go-resty/resty/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"

	"github.com/craftaholic/insider/internal/shared/config"
	"github.com/craftaholic/insider/internal/shared/constant"
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

	// Init resty client
	restyClient := resty.New()

	// Configure built-in retry
	restyClient.
		SetRetryCount(constant.RestMaxRetry).                        // Max 3 retries
		SetRetryWaitTime(constant.RestRetryWaitTime * time.Second).  // Initial wait
		SetRetryMaxWaitTime(constant.RestMaxWaitTime * time.Second). // Max wait time
		SetTimeout(constant.RestTimeOut * time.Second).              // Request timeout
		AddRetryCondition(func(r *resty.Response, err error) bool {
			// Retry on network errors
			if err != nil {
				return true
			}
			// Retry on specific HTTP status codes
			return r.StatusCode() >= 500 || r.StatusCode() == 429 // Server errors or rate limit
		}).
		SetRetryAfter(func(_ *resty.Client, resp *resty.Response) (time.Duration, error) {
			// Custom backoff - exponential with jitter
			retryCount := resp.Request.Attempt
			backoff := time.Duration(math.Pow(constant.RestExponentialBackOffScale, float64(retryCount))) * time.Second
			return backoff, nil
		})

	// Init Repository Layer
	messageRepository := repository.NewMessageRepository(db)
	cacheRepository := repository.NewCacheRepository(redisClient)
	notificationService := repository.NewNotificationService(
		restyClient,
		config.Env.WebhookAuthKey,
		config.Env.WebhookURL,
	)

	// Init Usecase Layer
	messageUsecase := usecase.NewMessageUsecase(
		messageRepository,
		cacheRepository,
		notificationService,
		config.Env.WorkerChanBuffer,
		config.Env.WorkerCount,
		config.Env.MessageCronDuration,
		config.Env.MessageBatchNumber,
	)

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
