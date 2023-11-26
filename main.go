package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/empfaze/golang_url_reducer/internal/config"
	"github.com/empfaze/golang_url_reducer/internal/http_server/handlers/redirect"
	"github.com/empfaze/golang_url_reducer/internal/http_server/handlers/url"
	"github.com/empfaze/golang_url_reducer/internal/logger"
	"github.com/empfaze/golang_url_reducer/internal/storage/sqlite"
	"github.com/empfaze/golang_url_reducer/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	config := config.MustLoad()
	logger := logger.SetupLogger(config.Env)

	storage, err := sqlite.New(config.StoragePath)
	if err != nil {
		logger.Error("Failed to init storage", utils.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url_reducer", utils.AllowedUsers()))
		r.Post("/", url.New(logger, storage))
	})

	router.Route("/alias", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url_reducer", utils.AllowedUsers()))
		r.Get("/{alias}", redirect.New(logger, storage))
	})

	logger.Info("Starting server", slog.String("address", config.Address))

	server := &http.Server{
		Addr:         config.Address,
		Handler:      router,
		ReadTimeout:  config.Timeout,
		WriteTimeout: config.Timeout,
		IdleTimeout:  config.IdleTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		logger.Error("Failed to start server")
	}

	logger.Error("Server has been stopped")
}
