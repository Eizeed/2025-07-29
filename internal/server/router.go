package server

import (
	"context"
	"net/http"
	"time"

	"github.com/Eizeed/2025-07-29/internal/pkg/config"
	"github.com/Eizeed/2025-07-29/internal/pkg/ctx"
	"github.com/Eizeed/2025-07-29/internal/server/handlers"
)

func initRoutes(mux *http.ServeMux, appCfg *config.AppConfig) {
	mux.HandleFunc("GET /api/v1/archive", wrapWithCfg(appCfg, handlers.GetArchiveList))
	mux.HandleFunc("POST /api/v1/archive", wrapWithCfg(appCfg, handlers.CreateArchive))
	mux.HandleFunc("GET /api/v1/archive/{zipName}", wrapWithCfg(appCfg, handlers.GetArchive))
	mux.HandleFunc("GET /api/v1/task", wrapWithCfg(appCfg, handlers.GetTasks))
	mux.HandleFunc("POST /api/v1/task", wrapWithCfg(appCfg, handlers.CreateTask))
	mux.HandleFunc("PATCH /api/v1/task/{uuid}", wrapWithCfg(appCfg, handlers.AddToTask))
	mux.HandleFunc("GET /api/v1/task/{uuid}", wrapWithCfg(appCfg, handlers.CheckTask))
}

func wrapWithCfg(appCfg *config.AppConfig, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().UTC()
		appCfg.Logger.Info("Reqest accepted. Path: ", r.URL, ", Method: ", r.Method)
		newCtx := context.WithValue(r.Context(), ctx.AppConfigKey{}, appCfg)
		r = r.WithContext(newCtx)

		handler(w, r)
		appCfg.Logger.Info("Reqest resolved. Time elapsed: ", time.Since(now).Milliseconds())
	}
}
