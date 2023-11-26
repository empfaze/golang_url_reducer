package url

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/empfaze/golang_url_reducer/internal/storage"
	"github.com/empfaze/golang_url_reducer/utils"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `jso:"alias,omitempty"`
}

type Response struct {
	utils.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave, alias string) (int64, error)
}

const OPERATION_TRACE_NEW = "internal.http_server.handlers.url.New"
const ALIAS_LENGTH = 5

func New(logger *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger = logger.With(
			slog.String("operation", OPERATION_TRACE_NEW),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var request Request

		err := render.DecodeJSON(r.Body, &request)
		if err != nil {
			logger.Error("Failed to decode request body", err)
			render.JSON(w, r, utils.Error("Failed to decode request"))

			return
		}

		logger.Info("Request body decoded", slog.Any("request", request))

		if err := validator.New().Struct(request); err != nil {
			validationErrors := err.(validator.ValidationErrors)

			logger.Error("Invalid request", err)
			render.JSON(w, r, utils.ValidateError(validationErrors))

			return
		}

		alias := request.Alias
		if alias == "" {
			alias = utils.NewRandomString(ALIAS_LENGTH)
		}

		id, err := urlSaver.SaveURL(request.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			logger.Info("Url already exists", slog.String("url", request.URL))
			render.JSON(w, r, utils.Error("Url already exists"))

			return
		}
		if err != nil {
			logger.Error("Failed to add url", err)
			render.JSON(w, r, utils.Error("Failed to add url"))

			return
		}

		logger.Info("Url has been added", slog.Int64("id", id))
		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{Response: utils.OK(), Alias: alias})
}
