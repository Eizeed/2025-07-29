package handlers

import (
	"encoding/json"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	stdIo "io"

	"github.com/Eizeed/2025-07-29/internal/pkg/archive"
	"github.com/Eizeed/2025-07-29/internal/pkg/io"
)

func CreateArchive(w http.ResponseWriter, r *http.Request) {
	type params struct {
		URLs []string `json:"urls"`
	}

	decoder := json.NewDecoder(r.Body)
	p := params{}

	decoder.Decode(&p)

	if len(p.URLs) > 3 {
		responseWithError(w, 400, "urls.len should be less or equal to 3")
		return
	}

	archive := archive.NewArchive()

	failed := []string{}
	succeeded := []string{}

	client := http.DefaultClient
	for _, url := range p.URLs {
		res, err := client.Get(url)
		if err != nil {
			continue
		}

		contentType := res.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "image/jpeg") {
			failed = append(failed, url)
			log.Println("Bad content type:", contentType)
			continue
		}

		blobName := res.Header.Get("Content-Disposition")
		if blobName == "" {
			failed = append(failed, url)
			log.Println("Bad content disposition:", blobName)
			continue
		}

		_, data, err := mime.ParseMediaType(blobName)
		if err != nil {
			failed = append(failed, url)
			log.Println("Error parsing media type:", err)
			continue
		}

		name := data["filename"]
		if name == "" {
			failed = append(failed, url)
			log.Println("Bad filename")
			continue
		}

		bytes, err := stdIo.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			failed = append(failed, url)
			log.Println("Failed to read from res.Body")
			continue
		}

		filePath, err := io.SaveToFileDir(name, bytes)
		if err != nil {
			failed = append(failed, url)
			log.Println("Failed to save", err)
			continue
		}

		archive.AddPath(filePath)
		succeeded = append(succeeded, url)
	}

	path, err := io.ZipFromArchive(&archive)
	if err != nil {
		log.Println("Failed to zip archive", err)
		responseWithError(w, 400, "Failed to zip archive "+err.Error())
	}

	type ResBody struct {
		Succeeded []string `json:"succeeded"`
		Failed    []string `json:"failed"`
		LocalPath string   `json:"local_path"`
		HttpPath  string   `json:"http_path"`
	}

	httpPath := "http://localhost:8080/api/v1/archive/" + filepath.Base(path)

	resBody := ResBody{
		Succeeded: succeeded,
		Failed:    failed,
		LocalPath: path,
		HttpPath:  httpPath,
	}

	responseWithBody(w, 201, resBody)
}

func GetArchive(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	// It's okay to take 4th because of router
	// This route has parts[0] which is ""
	// and after that we get "api", "v1" and our file name
	zipName := parts[4]
	zipDirPath, err := io.ZipDirPath()
	if err != nil {
		responseWithError(w, 400, "Failed to open zip dir path: "+err.Error())
		return
	}

	file, err := os.Open(filepath.Join(zipDirPath, zipName))
	if err != nil {
		responseWithError(w, 400, "Failed to open a file: "+err.Error())
		return
	}

	// Idk why it sets 200 status code implicitly
	_, err = stdIo.Copy(w, file)

	if err != nil {
		responseWithError(w, 400, "Failed to copy data from file: "+err.Error())
		return
	}
}
