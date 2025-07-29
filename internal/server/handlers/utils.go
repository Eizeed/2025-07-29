package handlers

import (
	"encoding/json"
	"net/http"
)

func responseWithBody(w http.ResponseWriter, statusCode int, body any) {
	bytes, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(bytes)
}

func responseWithError(w http.ResponseWriter, statusCode int, errMsg string) {
	type errorMessage struct {
		Error string `json:"error"`
	}

	err := errorMessage{
		Error: errMsg,
	}

	responseWithBody(w, statusCode, err)
}
