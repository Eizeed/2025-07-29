package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/Eizeed/2025-07-29/pkg/uuid"
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

func checkContentType(url string, client *http.Client) error {
	res, err := client.Head(url)
	if err != nil {
		return err
	}

	contentType := res.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/jpeg") && !strings.HasPrefix(contentType, "application/pdf") {
		return errors.New(fmt.Sprint("Bad content type: ", contentType))
	}

	return nil
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
	if !strings.HasPrefix(contentType, "image/jpeg") && !strings.HasPrefix(contentType, "application/pdf") {
		return FileRes{}, errors.New(fmt.Sprintln("Bad content type: ", contentType))
	}

	ext := ""
	switch contentType {
	case "image/jpeg":
		ext = ".jpeg"
	case "application/pdf":
		ext = ".pdf"
	default:
		return FileRes{}, errors.New("Invalid Content-Type")
	}

	blobName := res.Header.Get("Content-Disposition")
	name := uuid.NewV4().String() + ext
	if blobName != "" {
		_, data, err := mime.ParseMediaType(blobName)
		if err == nil {
			filename := data["filename"]
			if filename == "" {
				name = uuid.NewV4().String()
			} else {
				name = filename
			}
		}
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
