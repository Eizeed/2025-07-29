package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Eizeed/2025-07-29/internal/pkg/config"
	"github.com/Eizeed/2025-07-29/internal/pkg/task"
	"github.com/Eizeed/2025-07-29/pkg/dotenv"
)

func StartServer() {
	dotenv.DotEnv()

	mux := http.NewServeMux()

	appCfg := &config.AppConfig{
		TaskQueue: task.NewQueue(),
	}

	initRoutes(mux, appCfg)

	port := "8080"
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("Failed to listen and server: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Failed to Shutdown server: ", err)
	}

	println()
}
