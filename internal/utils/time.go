package utils

import (
	"time"

	"github.com/craftaholic/insider/internal/shared/config"
)

func GetContextTimeout() time.Duration {
	return time.Duration(config.Env.ContextTimeout) * time.Second
}
