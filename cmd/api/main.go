package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/genda/genda-api/cmd/api/internal"
	"github.com/genda/genda-api/internal/storage/postgres"
	"github.com/genda/genda-api/pkg/config"
	"github.com/pkg/errors"
	"github.com/rs/cors"
)

var buildVersion string

func run() error {

	log := log.New(os.Stdout, "[cmd.api.main.run] ", log.LstdFlags|log.Lmicroseconds|log.Lmsgprefix)

	conf := config.New()

	var postgresDB postgres.PostgresDB
	postgresDB, _ = postgres.NewConnection()

	//var rdb redis.Redis
	//rdb, _ = redis.NewConnection()

	log.Printf("Starting genda API in %s mode ...", conf.Environment)

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://*.genda.com.br", "https://*.genda.io", "https://calender-dev.genda.group", "https://calender.genda.group", "http://localhost:3000"},
		AllowedMethods:   []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	api := http.Server{
		Addr:    ":" + conf.APIPort,
		Handler: c.Handler(internal.API(conf.Environment, shutdown, log, postgresDB)), //add redis again if we're to use it
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("build version %s", buildVersion)
		log.Printf("listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "Internal Server Error")

	case sig := <-shutdown:
		log.Printf("%v : Starting Shutdown", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Second) //conf.API.ShutdownTimeout
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("Graceful Shutdown wasnt complete in %v : %v", 10, err)
			err = api.Close()
		}

		// Log the status of this shutdown.
		switch {
		case sig == syscall.SIGSTOP:
			return errors.New("A fail caused shutdown")
		case err != nil:
			return errors.Wrap(err, "could not stop gracefully")
		}
	}

	return nil

}

func main() {

	if err := run(); err != nil {
		log.Println("error:", err)
		os.Exit(1)
	}

}
