package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"

	stdIo "io"

	"github.com/Eizeed/2025-07-29/internal/pkg/archive"
	"github.com/Eizeed/2025-07-29/internal/pkg/constants"
	"github.com/Eizeed/2025-07-29/internal/pkg/io"
)

func GetArchiveList(w http.ResponseWriter, r *http.Request) {
	zipDirPath, err := io.ZipDirPath()
	if err != nil {
		responseWithError(w, 400, "Failed to open zip dir path: "+err.Error())
		return
	}

	entries, err := os.ReadDir(zipDirPath)
	if err != nil {
		responseWithError(w, 400, "Failed to read zip dir: "+err.Error())
		return
	}

	type zipPath struct {
		LocalPath string `json:"local_path"`
		HttpPath  string `json:"http_path"`
	}

	res := []zipPath{}

	for _, entry := range entries {
		ext := path.Ext(entry.Name())
		if ext == ".zip" {
			localPath := filepath.Join(zipDirPath, entry.Name())
			httpPath := "http://localhost:8080/api/v1/archive/" + filepath.Base(entry.Name())

			res = append(res, zipPath{
				LocalPath: localPath,
				HttpPath:  httpPath,
			})
		}
	}

	type ResBody struct {
		Paths []zipPath `json:"paths"`
	}

	responseWithBody(w, 200, ResBody{
		Paths: res,
	})
}

func CreateArchive(w http.ResponseWriter, r *http.Request) {
	type params struct {
		URLs []string `json:"urls"`
	}

	decoder := json.NewDecoder(r.Body)
	p := params{}

	err := decoder.Decode(&p)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintln("Failed to decode body: ", err))
		return
	}

	if len(p.URLs) > constants.URL_LIMIT {
		responseWithError(w, 400, fmt.Sprint("urls.len should be less or equal to ", constants.URL_LIMIT))
		return
	}

	client := http.DefaultClient

	for _, u := range p.URLs {
		err = checkContentType(u, client)
		if err != nil {
			responseWithError(w, 400, err.Error())
			return
		}
	}

	archive := archive.NewArchive()

	failed := []string{}
	succeeded := []string{}

	type parseRes struct {
		fileRes FileRes
		url     string
		err     error
	}
	resCh := make(chan parseRes, len(p.URLs))

	wg := sync.WaitGroup{}
	wg.Add(len(p.URLs))
	for _, url := range p.URLs {
		go func(url string) {
			defer wg.Done()

			fileRes, err := parseFile(url, client)
			if err != nil {
				log.Println(err.Error())
				resCh <- parseRes{
					fileRes: FileRes{},
					url:     url,
					err:     err,
				}
				return
			}

			resCh <- parseRes{
				fileRes: fileRes,
				url:     url,
				err:     nil,
			}
		}(url)
	}

	wg.Wait()
	close(resCh)

	for res := range resCh {
		if res.err != nil {
			failed = append(failed, res.url)
			log.Println(res.err.Error())
			continue
		}

		filePath, err := io.SaveToFileDir(res.fileRes.name, res.fileRes.bytes)
		if err != nil {
			failed = append(failed, res.url)
			log.Println(res.err.Error())
			continue
		}

		succeeded = append(succeeded, res.url)

		// Can't push more than archive can hold
		// because we check amount of urls at start
		// and it is oneshot operation
		_ = archive.AddPath(filePath)
	}

	path, err := io.ZipFromArchive(&archive)
	if err != nil {
		log.Println("Failed to zip archive", err)
		responseWithError(w, 400, "Failed to zip archive "+err.Error())
		return
	}
	httpPath := "http://localhost:8080/api/v1/archive/" + filepath.Base(path)

	type ResBody struct {
		Succeeded []string `json:"succeeded"`
		Failed    []string `json:"failed"`
		LocalPath string   `json:"local_path"`
		HttpPath  string   `json:"http_path"`
	}

	resBody := ResBody{
		Succeeded: succeeded,
		Failed:    failed,
		LocalPath: path,
		HttpPath:  httpPath,
	}

	responseWithBody(w, 201, resBody)
}

func GetArchive(w http.ResponseWriter, r *http.Request) {
	zipName := r.PathValue("zipName")

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
