package main

import (
	"context"
	"github.com/Sambitmohanty954/students-api-golang/internal/config"
	"github.com/Sambitmohanty954/students-api-golang/internal/http/handlers/student"
	"github.com/Sambitmohanty954/students-api-golang/internal/storage/sqlite"
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
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("Storage initialized", slog.String(" env path", cfg.Env), slog.String("version", "1.0.0"))

	// setup Router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))

	// Setup server
	httpServer := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	// PrintF we use to concatination,
	slog.Info("Server started ", slog.String("address ", cfg.Address))

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

	err = httpServer.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown http server", slog.String("error", err.Error()))
	}

	slog.Info("Server shut down successfully")
}
