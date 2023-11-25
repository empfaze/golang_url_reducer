package main

import (
	"os"

	"github.com/empfaze/golang_url_reducer/internal/config"
	"github.com/empfaze/golang_url_reducer/internal/logger"
	"github.com/empfaze/golang_url_reducer/internal/storage/sqlite"
	"github.com/empfaze/golang_url_reducer/utils"

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

	_ = storage
}
