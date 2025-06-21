package config

import (
	"os"
	"strconv"

	"github.com/craftaholic/insider/internal/shared/constant"
	"github.com/craftaholic/insider/internal/shared/log"
	"github.com/joho/godotenv"
)

var logger log.Log

var Env *EnvConfig

type EnvConfig struct {
	// App config
	AppEnv         string
	ContextTimeout int
	ServerAddress  string

	// DB config
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSslMode  string

	// Redis config
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       string

	// Notification service
	WebhookURL     string
	WebhookAuthKey string
	WebhookTimeout int

	// Concurency config
	MessageBatchNumber  int
	MessageCronDuration int
}

func LoadEnv() {
	logger = log.BaseLogger.WithFields("bootstrap", "LoadEnv")

	// Load .env file first (won't override real env vars)
	_ = godotenv.Load(".env")

	env := &EnvConfig{
		AppEnv:         getEnv("APP_ENV", "development"),
		ContextTimeout: getIntEnv("CONTEXT_TIMEOUT", constant.DefaultContextTimeOut),
		ServerAddress:  getEnv("SERVER_ADDR", "8080"),

		// DB config
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "message_system"),
		DBUser:     getEnv("DB_USER", "posgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres123"),
		DBSslMode:  getEnv("DB_SSL_MODE", "disable"),

		// Redis config
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnv("REDIS_DB", "0"),

		// Notification service
		WebhookURL:     getEnvOrPanic("WEBHOOK_URL"),
		WebhookAuthKey: getEnvOrPanic("WEBHOOK_URL"),
		WebhookTimeout: getIntEnv("WEBHOOK_TIMEOUT", 30),

		// Concurency config
		MessageBatchNumber:  getIntEnv("MESSAGE_BATCH_NUMBER", 2),
		MessageCronDuration: getIntEnv("MESSAGE_BATCH_NUMBER", 120),
	}

	logger.Info("Loaded Config", "AppEnv", env.AppEnv)
	Env = env
}

// getEnv gets an environment variable or returns a default value.
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getIntEnv(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

// getEnvOrPanic gets an environment variable or panics if not set.
func getEnvOrPanic(key string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	logger.Fatal("Required environment variable missing", "key", key)
	return ""
}
