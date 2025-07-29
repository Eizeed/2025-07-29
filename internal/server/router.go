package server

import (
	"net/http"

	"github.com/Eizeed/2025-07-29/internal/server/handlers"
)

func initRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/archive", post(handlers.CreateArchive))
	mux.HandleFunc("/api/v1/archive/", get(handlers.GetArchiveStatus))
	// mux.HandleFunc("/task", handler)
	// mux.HandleFunc("/task/{index}", handler)
	// mux.HandleFunc("/task/{index}", handler)
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
		println("METHOD:", r.Method)
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
