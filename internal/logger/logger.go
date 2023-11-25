package logger

import (
	"log/slog"
	"os"

	"github.com/empfaze/golang_url_reducer/utils"
)

func SetupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case utils.ENV_LOCAL:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case utils.ENV_DEV:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case utils.ENV_PROD:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return logger
}
