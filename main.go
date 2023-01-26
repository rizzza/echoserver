package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rizzza/echoserver/internal/config"
	"github.com/rizzza/echoserver/internal/handlers"
	log "github.com/sirupsen/logrus"
	"go.uber.org/automaxprocs/maxprocs"
)

var (
	Version = "dev"
)

func main() {
	undo, err := maxprocs.Set(maxprocs.Logger(log.Printf))
	if err != nil {
		log.Fatalf("failed to set GOMAXPROCS: %v", err)
	}

	defer undo()

	cfg := config.Get()
	lvl, _ := log.ParseLevel(cfg.LogLevel)
	log.SetLevel(lvl)

	srv := http.Server{
		Handler:      handlers.New(cfg).GetRouter(),
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.WriteTimeout)
		defer cancel()

		_ = srv.Shutdown(ctx)
	}()

	go func() {
		addr := fmt.Sprintf("0.0.0.0:%s", cfg.Port)
		log.Printf("Listening on %s", addr)
		if !errors.Is(srv.ListenAndServe(), http.ErrServerClosed) {
			log.Fatal("failed to start up server", err)
		}
	}()

	<-shutdown
	log.Println("Shutting down server...")
}
