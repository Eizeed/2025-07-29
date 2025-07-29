package server

import "net/http"

func StartServer() {
	mux := http.NewServeMux()

	initRoutes(mux)

	port := "8080"

	http.ListenAndServe(":"+port, mux)
}
