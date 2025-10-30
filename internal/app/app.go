package app

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.opencensus.io/trace"
)

// contextID is the type that we use to context keys.
type contextID struct{}

// ContextIDKey is the the identifier that we use to get values from context
var ContextIDKey = &contextID{}

// Handler is type that handles an http request in our application
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, params httprouter.Params) error

// CtxValue type holds all info that we want pass to context.
// It holds: TraceID, HTTPStatusCode and Now(that is the current time)
type CtxValue struct {
	TraceID        string
	HTTPStatusCode int
	Now            time.Time
}

type App struct {
	*httprouter.Router
	shutdown chan os.Signal
	mw       []Middleware
}

// New creates an App value that handle a set of routes for the application.
func New(shutdown chan os.Signal, mw ...Middleware) *App {
	a := App{
		Router:   httprouter.New(),
		shutdown: shutdown,
		mw:       mw,
	}

	return &a
}

// SignalShutdown gracefully shutdown the app if something happens wrong
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// Handle is the function for mounting Handlers given a HTTP verb and path pair.
func (a *App) Handle(verb, path string, handler Handler, mw ...Middleware) {

	// Add specific middlewares around this handler.
	handler = addMiddleware(mw, handler)

	// Add the app default middlewares to the chain.
	handler = addMiddleware(a.mw, handler)

	// This function will run on each request
	h := func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx, span := trace.StartSpan(r.Context(), "internal.app")
		defer span.End()

		v := CtxValue{
			TraceID: span.SpanContext().TraceID.String(),
			Now:     time.Now(),
		}
		ctx = context.WithValue(ctx, ContextIDKey, &v)

		// Call the wrapped handler functions.
		if err := handler(ctx, w, r, params); err != nil {
			a.SignalShutdown()
			return
		}
	}

	// Add this handler to the route.
	a.Router.Handle(verb, path, h)
}

// It overrides the ServeHTTP of the embedded httprouter.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Router.ServeHTTP(w, r)
}
