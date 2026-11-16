package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

// dependency injection, contain dependencies for handlers, helpers and middlewares
type application struct {
	config config
	logger *slog.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "server port number")
	flag.StringVar(&cfg.env, "environment", "development", "(development|staging|production)")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdin, nil))

	app := &application{
		config: cfg,
		logger: logger,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("Starting Server ", "addr", cfg.port, "env", cfg.env)
	err := server.ListenAndServe()

	logger.Error(err.Error())
	os.Exit(1)
}
