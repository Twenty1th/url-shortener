package save

import (
	"UrlShortner/internal/lib/api/response"
	"UrlShortner/internal/lib/random"
	"UrlShortner/internal/storage"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

const aliasLength = 6

type Request struct {
	Url   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveUrl(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, saver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.save.url"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", "error", err)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}
		log.Info("request body decoded", slog.Any("request", req))
		if err := validator.New().Struct(req); err != nil {
			var validatorErr validator.ValidationErrors
			errors.As(err, &validatorErr)

			log.Error("failed to validate request", "error", err)
			render.JSON(w, r, response.ValidationError(validatorErr))
			return
		}
		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}
		id, err := saver.SaveUrl(req.Url, alias)
		if errors.Is(err, storage.ErrUrlExists) {
			log.Info("url already exists", slog.String("url", req.Url))
			render.JSON(w, r, response.Error("url already exists"))
			return

		}
		if err != nil {
			log.Error("failed to save url", "error", err)
			render.JSON(w, r, response.Error("failed to save url"))
			return
		}
		log.Info("url saved", slog.String("url", req.Url), slog.Int64("id", id))
		render.JSON(w, r, Response{
			Response: response.OK(),
			Alias:    alias,
		})
	}
}
