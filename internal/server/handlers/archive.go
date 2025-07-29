package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func CreateArchive(w http.ResponseWriter, r *http.Request) {
	log.Println("Create Archive")
	w.WriteHeader(200)
}

func GetArchiveStatus(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	uuidStr := parts[4]
	_, err := uuid.Parse(uuidStr)
	if err != nil {
		http.Error(w, "Invalid uuid format", 400)
		return
	}

	w.WriteHeader(200)
}
