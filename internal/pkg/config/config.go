package config

import (
	"github.com/Eizeed/2025-07-29/internal/pkg/log"
	"github.com/Eizeed/2025-07-29/internal/pkg/task"
)

type AppConfig struct {
	TaskQueue task.TaskQueue
	Logger    log.Logger
}
