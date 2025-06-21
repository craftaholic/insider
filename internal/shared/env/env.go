package env

import (
	"os"
	"strconv"

	"github.com/craftaholic/insider/internal/shared/log"
	"github.com/joho/godotenv"
)

var logger log.Log

var Env *EnvConfig

type EnvConfig struct {
	AppEnv         string
	ContextTimeout int
	ServerAddress  string
}

func LoadEnv() {
	logger = log.BaseLogger.WithFields("bootstrap", "LoadEnv")

	// Load .env file first (won't override real env vars)
	_ = godotenv.Load(".env")

	env := &EnvConfig{
		AppEnv:         getEnv("APP_ENV", "development"),
		ContextTimeout: getIntEnv("CONTEXT_TIMEOUT", 30),
		ServerAddress:  getEnv("SERVER_ADDR", "8080"),
	}

	logger.Info("Loaded Config", "AppEnv", env.AppEnv)
	Env = env
}

// getEnv gets an environment variable or returns a default value
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

// getEnvOrPanic gets an environment variable or panics if not set
func getEnvOrPanic(key string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	logger.Fatal("Required environment variable missing", "key", key)
	return ""
}
