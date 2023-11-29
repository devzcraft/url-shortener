package redirect

import (
	"context"
	"errors"
	"net/http"

	"log/slog"

	resp "github.com/devzcraft/url-shortener/internal/lib/api/response"
	"github.com/devzcraft/url-shortener/internal/lib/logger/sl"
	"github.com/devzcraft/url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// TODO: add mock
type URLGetter interface {
	GetURL(ctx context.Context, alias string) (string, error)
}

// TODO: add test
func New(ctx context.Context, log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		resURL, err := urlGetter.GetURL(ctx, alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, resp.Error("not found"))

			return
		}

		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("got url", slog.String("url", resURL))

		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
