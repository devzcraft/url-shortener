package remove

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	resp "github.com/devzcraft/url-shortener/internal/lib/api/response"
	"github.com/devzcraft/url-shortener/internal/lib/logger/sl"
	"github.com/devzcraft/url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// TODO: add mock
type URLRemover interface {
	DeleteURL(ctx context.Context, alias string) error
}

type Response struct {
	resp.Response
}

// TODO: add test
func New(ctx context.Context, log *slog.Logger, urlRemover URLRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.remove.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("not found"))

			return
		}

		err := urlRemover.DeleteURL(ctx, alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Error("failed to delete url", sl.Err(err))

			render.JSON(w, r, resp.Error("not found"))

			return
		}

		if err != nil {
			log.Error("failed to delete url", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}
}
