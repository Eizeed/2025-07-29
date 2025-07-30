package server

import (
	"context"
	"net/http"

	"github.com/Eizeed/2025-07-29/internal/pkg/config"
	"github.com/Eizeed/2025-07-29/internal/pkg/ctx"
	"github.com/Eizeed/2025-07-29/internal/server/handlers"
)

func initRoutes(mux *http.ServeMux, appCfg *config.AppConfig) {
	mux.HandleFunc("POST /api/v1/archive", wrapWithCfg(appCfg, handlers.CreateArchive))
	mux.HandleFunc("GET /api/v1/archive/{zipName}", wrapWithCfg(appCfg, handlers.GetArchive))
	mux.HandleFunc("POST /api/v1/task", wrapWithCfg(appCfg, handlers.CreateTask))
	mux.HandleFunc("GET /api/v1/task/completed", wrapWithCfg(appCfg, handlers.GetCompletedTasks))
	mux.HandleFunc("PATCH /api/v1/task/{uuid}", wrapWithCfg(appCfg, handlers.AddToTask))
	mux.HandleFunc("GET /api/v1/task/{uuid}", wrapWithCfg(appCfg, handlers.CheckTask))
}

func wrapWithCfg(appCfg *config.AppConfig, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newCtx := context.WithValue(r.Context(), ctx.AppConfigKey{}, appCfg)
		r = r.WithContext(newCtx)

		handler(w, r)
	}
}
