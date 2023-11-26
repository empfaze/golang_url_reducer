package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/empfaze/golang_url_reducer/internal/storage"
	"github.com/empfaze/golang_url_reducer/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

const OPERATION_TRACE_NEW = "internal.http_server.handlers.redirect.New"

func New(logger *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.With(
			slog.String("operation", OPERATION_TRACE_NEW),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			logger.Info("Alias is empty")
			render.JSON(w, r, utils.Error("Invalid request"))

			return
		}

		responseURL, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			logger.Info("Url not found", "alias", alias)
			render.JSON(w, r, utils.Error("Not found"))

			return
		}
		if err != nil {
			logger.Error("Failed to get url", err)
			render.JSON(w, r, utils.Error("Internal error"))

			return
		}

		logger.Info("Got url", slog.String("url", responseURL))

		http.Redirect(w, r, responseURL, http.StatusFound)
	}
}
