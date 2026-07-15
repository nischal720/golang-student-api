package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nischal720/golang-student-api/internal/config"
)

func main() {
	// Load config, database setup, router setup, server setup

	//load config

	cfg := config.MustLoad()

	//Database setup

	//router setup
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to student api"))
	})

	//server setup

	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}
	fmt.Println("Server Start")
	slog.Info("server start", slog.String("address", cfg.Address))
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()

		if err != nil {
			log.Fatal("Failed to start server")
		}
	}()

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	slog.Info("Shutting down the server")

	err := server.Shutdown(ctx)

	if err != nil {
		slog.Error("failed to down server", slog.String("error", err.Error()))
	}

	slog.Info("server shut down successfull")

}
