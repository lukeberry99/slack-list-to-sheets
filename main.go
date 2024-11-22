package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lukeberry99/slack-list-to-sheets/config"
	"github.com/lukeberry99/slack-list-to-sheets/handlers"
	"github.com/lukeberry99/slack-list-to-sheets/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create handlers
	csvHandler := handlers.NewCSVHandler(cfg.SlackToken)

	// Setup router
	mux := http.NewServeMux()
	mux.HandleFunc("/get-file", csvHandler.HandleGetCSV)

	// Add middleware
	handler := middleware.Logging(mux)

	// Create server
	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: handler,
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
		serverStopCtx()
	}()

	// Run the server
	log.Printf("Server started on port %s", cfg.ServerPort)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %v", err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
