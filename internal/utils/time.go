package utils

import (
	"time"

	"github.com/craftaholic/insider/internal/shared/env"
)

func GetContextTimeout() time.Duration {
	return time.Duration(env.Env.ContextTimeout) * time.Second
}
