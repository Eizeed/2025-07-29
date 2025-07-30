package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/Eizeed/2025-07-29/internal/pkg/constants"
	"github.com/Eizeed/2025-07-29/internal/pkg/ctx"
	"github.com/Eizeed/2025-07-29/internal/pkg/io"
	"github.com/Eizeed/2025-07-29/pkg/uuid"
)

func GetTasks(w http.ResponseWriter, r *http.Request) {
	appCfg := ctx.GetAppConfig(r.Context())

	type taskView struct {
		UUID    uuid.UUID `json:"uuid"`
		Content []string  `json:"content"`
	}

	type ResBody struct {
		Tasks []taskView `json:"tasks"`
	}

	tasks := []taskView{}

	for _, task := range appCfg.TaskQueue.ViewTasks() {
		tasks = append(tasks, taskView{
			UUID:    task.UUID,
			Content: task.Archive.Content,
		})
	}

	responseWithBody(w, 200, ResBody{
		Tasks: tasks,
	})
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	appCfg := ctx.GetAppConfig(r.Context())
	taskUuid, err := appCfg.TaskQueue.InsertTask()
	if err != nil {
		responseWithError(w, 400, err.Error())
		return
	}

	type ResBody struct {
		Uuid uuid.UUID `json:"uuid"`
	}

	responseWithBody(w, 201, ResBody{
		Uuid: taskUuid,
	})
}

func AddToTask(w http.ResponseWriter, r *http.Request) {
	uuidStr := r.PathValue("uuid")
	taskUuid, err := uuid.Parse(uuidStr)
	if err != nil {
		responseWithError(w, 400, err.Error())
		return
	}

	type params struct {
		URLs []string `json:"urls"`
	}

	decoder := json.NewDecoder(r.Body)
	p := params{}

	err = decoder.Decode(&p)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintln("Failed to decode body: ", err))
		return
	}

	if len(p.URLs) > constants.URL_LIMIT {
		responseWithError(w, 400, fmt.Sprint("urls.len should be less or equal to ", constants.URL_LIMIT))
		return
	}

	appCfg := ctx.GetAppConfig(r.Context())
	task, exists := appCfg.TaskQueue.GetTask(taskUuid)
	if !exists {
		responseWithError(w, 404, "Task not found")
		return
	}

	contentLen := len(task.Archive.Content)

	if contentLen == constants.URL_LIMIT {
		responseWithError(w, 400, fmt.Sprint("Archive already contains", constants.URL_LIMIT, "URLs"))
		return
	}

	if contentLen+len(p.URLs) > constants.URL_LIMIT {
		responseWithError(w, 400, fmt.Sprint("Archive contains ", contentLen, " out of", constants.URL_LIMIT, ". You provided ", len(p.URLs), " URLs"))
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
		_ = task.Archive.AddPath(filePath)
	}

	type ResBody struct {
		TaskUUID  uuid.UUID `json:"task_uuid"`
		Succeeded []string  `json:"succeeded"`
		Failed    []string  `json:"failed"`
	}

	resBody := ResBody{
		TaskUUID:  task.UUID,
		Succeeded: succeeded,
		Failed:    failed,
	}

	responseWithBody(w, 201, resBody)
}

func CheckTask(w http.ResponseWriter, r *http.Request) {
	uuidStr := r.PathValue("uuid")
	taskUuid, err := uuid.Parse(uuidStr)
	if err != nil {
		responseWithError(w, 400, err.Error())
		return
	}

	appCfg := ctx.GetAppConfig(r.Context())

	task, exists := appCfg.TaskQueue.GetTask(taskUuid)
	if !exists {
		responseWithError(w, 404, "Task not found")
		return
	}

	contentLen := len(task.Archive.Content)
	statusStr := fmt.Sprint(contentLen, " out of ", constants.URL_LIMIT)

	localPath := ""
	httpPath := ""
	if contentLen == constants.URL_LIMIT {
		localPath, err = io.ZipFromArchive(&task.Archive)
		if err != nil {
			log.Println("Failed to zip archive", err)
			responseWithError(w, 400, "Failed to zip archive "+err.Error())
			return
		}
		httpPath = "http://localhost:8080/api/v1/archive/" + filepath.Base(localPath)

		appCfg.TaskQueue.RemoveByUUID(task.UUID)
	}

	type ResBody struct {
		TaskUUID         uuid.UUID `json:"task_uuid"`
		ArchiveUUID      uuid.UUID `json:"archive_uuid"`
		ArchuiveContent  []string  `json:"archive_content"`
		Status           string    `json:"status"`
		ArchiveLocalPath string    `json:"archive_local_path,omitempty"`
		ArchiveHttpPath  string    `json:"archive_http_path,omitempty"`
	}

	responseWithBody(w, 200, ResBody{
		TaskUUID:         task.UUID,
		ArchiveUUID:      task.Archive.UUID,
		ArchuiveContent:  task.Archive.Content,
		Status:           statusStr,
		ArchiveLocalPath: localPath,
		ArchiveHttpPath:  httpPath,
	})
}
