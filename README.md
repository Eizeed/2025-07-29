## Routes
#### GET http://localhost:8080/api/v1/archive

- req: None <br/>
- res: ```{ "paths": ["local_path": string, "http_path": string] }``` <br/>
- exp: Returns every archive in zipDirPath<br/>

#### POST http://localhost:8080/api/v1/archive

- req: ```{ "urls": [string] }``` <br/>
- res: ```{ "succeeded": [string], "failed": [string], "local_path": string, "http_path": string }``` <br/>
- exp: Downloads files, zip them and return to user<br/>

#### GET http://localhost:8080/api/v1/archive/{zipName}

- params: ```zipName: {zipUuid}.zip``` <br/>
- req: None <br/>
- res: binary ZIP-archive <br/>
- exp: Tries to get zip-archive with given name<br/>

#### GET http://localhost:8080/api/v1/task

- req: None <br/>
- res: ```{ "tasks": [ { "uuid": string, "content": [string] } ] }``` <br/>
- exp: Gets every current task in queue<br/>

#### POST http://localhost:8080/api/v1/task

- req: None <br/>
- res: ```{ "uuid": string }``` <br/>
- exp: Creates a task. Will fail if there are 3 tasks in queue<br/>

#### PATCH http://localhost:8080/api/v1/task/{uuid}

- params: ```uuid: string``` <br/>
- req: ```{ "urls": [string] }``` <br/>
- res: ```{ "task_uuid": string, "succeeded": [string], "failed": [string] }``` <br/>
- exp: Downloads files and saves them in tasks<br/>

#### GET http://localhost:8080/api/v1/task/{uuid}

- params: ```uuid: string```<br/>
- req: None <br/>
- res: if archive is full ```{ "task_uuid": string, "archive_uuid": string, "archive_content": [string], "status": string, "archive_local_path": string, "archive_http_path": string }```<br/>
- res: if archive is not full ```{ "task_uuid": string, "archive_uuid": string, "archive_content": [string], "status": string }```<br/>
- exp: Show status of Task. If task has 3 files it will return 2 urls to zip archive.<br/>
- Task will be removed from queue and user will be able to get archive with ```GET http://localhost:8080/api/v1/archive/{zipName}```<br/>
Or find its paths with ```GET http://localhost:8080/api/v1/archive```<br/>

## ENV
<p>
There are 4 env variables configured in app:<br/>
</p>

- PORT - port for server to listen on. Default is `8080`
- ZIP_PATH - output dir for zip files. Default is `$PWD/zip/`
- FILE_PATH - output dir for downloaded files. Default is `$PWD/files/`
- LOG_LEVEL - possible values: 0 - NONE, 1 - ERROR, 2 - WARN, 3 - ERROR, 4 - DEBUG. Default is DEBUG

## About project
I thought this project is not hard. So even though I'm not fluent in go, i decided to not use ANY dependency besides std.<br/>
That was my challenge to make something working without relying on 3rd Party libs like `chi`, `uuid`, `dotenv`<br/>
Also it was really interesting for me to dive into those things as I didn't do something simmilar before<br/>
My implementaitons might not be ideal but they work on that scale.

## Notes
- Didn't use `cmp` directory because I only need 1 binary. Also the project is not so big.
- Used `goroutines` while parsing files
- `RWMutex` is used on `TaskQueue` might need to read more than write
- `Task` has `mutex`. With `mutex` for each `Task` we can block `Task` and not block `TaskQueue`
which will be useful in case where multiple users make a change request but on different `Task`s
- Restriting types and amount of url on request. If they are invalid request will return an error
- `Assert` for `AppConfig` from `context` because if there is no config, it means there is a bug in program
- If downloaded file has no `Content-Disposition` and `filename` it will be replace with `uuid`+ ext from `Content-Type`
- If response has no `Content-Type` server will return error as the server should only process `jpeg` and `png`
