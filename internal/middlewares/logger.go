package middlewares

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/genda/genda-api/internal/app"
	"github.com/julienschmidt/httprouter"
	"go.opencensus.io/trace"
)

func Logger(log *log.Logger) app.Middleware {

	currentMw := func(current app.Handler) app.Handler {

		handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
			ctx, span := trace.StartSpan(ctx, "internal.middlewares.logger")
			defer span.End()

			// If the context is missing this value, request the service
			// to be shutdown gracefully.
			v, ok := ctx.Value(app.ContextIDKey).(*app.CtxValue)
			if !ok {
				return app.NewShutdownError("app value doesnt exists in context")
			}

			err := current(ctx, w, r, params)

			log.Printf("%s : (%d) : %s %s - %s (%s)",
				v.TraceID, v.HTTPStatusCode, r.Method, r.URL.Path, r.RemoteAddr, time.Since(v.Now),
			)

			return err
		}

		return handler
	}

	return currentMw
}
