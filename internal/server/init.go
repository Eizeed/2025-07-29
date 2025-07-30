package server

import (
	"net/http"

	"github.com/Eizeed/2025-07-29/internal/pkg/task"
	"github.com/Eizeed/2025-07-29/internal/pkg/config"
)

func StartServer() {
	mux := http.NewServeMux()

	appCfg := &config.AppConfig {
		TaskQueue: task.NewQueue(),
	}

	initRoutes(mux, appCfg)

	port := "8080"

	http.ListenAndServe(":"+port, mux)
}
