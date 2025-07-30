package server

import (
	"context"
	stdLog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Eizeed/2025-07-29/internal/pkg/config"
	"github.com/Eizeed/2025-07-29/internal/pkg/log"
	"github.com/Eizeed/2025-07-29/internal/pkg/task"
	"github.com/Eizeed/2025-07-29/pkg/dotenv"
)

func StartServer() {
	stdLog.Println("Parsing .env file...")
	err := dotenv.DotEnv()
	if err != nil {
		stdLog.Println(err)
	} else {
		stdLog.Println(".env file is parsed")
	}

	mux := http.NewServeMux()

	logLevel := os.Getenv("LOG_LEVEL")
	level, err := log.LogLevelFromStr(logLevel)
	if err != nil {
		level = log.DEBUG
	}

	appCfg := &config.AppConfig{
		TaskQueue: task.NewQueue(),
		Logger:    log.NewLogger(level),
	}

	initRoutes(mux, appCfg)
	stdLog.Println("Router initialized")

	stdLog.Println("Setting port...")
	port := "8080"
	envPort := os.Getenv("PORT")

	if envPort != "" {
		port = envPort
	}
	stdLog.Println("Port set to ", port)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		stdLog.Println("Server starts listening...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			stdLog.Fatalln("Failed to listen and server: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	stdLog.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		stdLog.Fatal("Failed to Shutdown server: ", err)
	}

	stdLog.Println("Server exited")

	println()
}
