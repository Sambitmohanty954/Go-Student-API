package main

import (
	"context"
	"github.com/Sambitmohanty954/students-api-golang/internal/config"
	"github.com/Sambitmohanty954/students-api-golang/internal/http/handlers/student"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// go run cmd/student-api/main.go -config config/local.yaml  (to run this project)
func main() {
	// Load config
	cfg := config.MustLoad()
	// Database setup

	// setup Router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New())

	// Setup server
	httpServer := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	// PrintF we use to concatination,
	slog.Info("Server started ", slog.String("address ", cfg.HTTPServer.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start http server")
		}

	}()
	<-done

	slog.Info("Shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := httpServer.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown http server", slog.String("error", err.Error()))
	}

	slog.Info("Server shut down successfully")
}
