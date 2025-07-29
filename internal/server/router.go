package server

import (
	"context"
	"net/http"

	"github.com/Eizeed/2025-07-29/internal/pkg/config"
	"github.com/Eizeed/2025-07-29/internal/pkg/ctx"
	"github.com/Eizeed/2025-07-29/internal/server/handlers"
)

func initRoutes(mux *http.ServeMux, appCfg *config.AppConfig) {
	mux.HandleFunc("/api/v1/archive", wrapWithCfg(appCfg, post(handlers.CreateArchive)))
	mux.HandleFunc("/api/v1/archive/", wrapWithCfg(appCfg, get(handlers.GetArchive)))
	// mux.HandleFunc("/task", wrapWithCfg(get(handler)))
	// mux.HandleFunc("/task/{index}", wrapWithCfg(patch(handler)))
	// mux.HandleFunc("/task/{index}", wrapWithCfg(get(handler)))
}

func wrapWithCfg(appCfg *config.AppConfig, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newCtx := context.WithValue(r.Context(), ctx.AppConfigKey{}, appCfg)
		r = r.WithContext(newCtx)

		handler(w, r)
	}
}

func get(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		handler(w, r)
	}
}

func post(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		handler(w, r)
	}
}

func patch(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		handler(w, r)
	}
}
