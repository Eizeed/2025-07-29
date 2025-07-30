package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"
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

type FileRes struct {
	name  string
	bytes []byte
}

func parseFile(url string, client *http.Client) (FileRes, error) {
	res, err := client.Get(url)
	if err != nil {
		return FileRes{}, err
	}

	contentType := res.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/jpeg") {
		return FileRes{}, errors.New(fmt.Sprintln("Bad content type:", contentType))
	}

	blobName := res.Header.Get("Content-Disposition")
	if blobName == "" {
		return FileRes{}, errors.New(fmt.Sprintln("Bad content disposition:", blobName))
	}

	_, data, err := mime.ParseMediaType(blobName)
	if err != nil {
		return FileRes{}, errors.New(fmt.Sprintln("Error parsing media type:", err))
	}

	name := data["filename"]
	if name == "" {
		return FileRes{}, errors.New(fmt.Sprintln("Bad filename"))
	}

	bytes, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return FileRes{}, errors.New(fmt.Sprintln("Failed to read from res.Body"))
	}

	return FileRes{
		name:  name,
		bytes: bytes,
	}, nil
}
