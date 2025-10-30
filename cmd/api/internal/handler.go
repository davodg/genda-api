package internal

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/genda/genda-api/cmd/api/internal/routes"
	routes "github.com/genda/genda-api/cmd/api/internal/routes"

	"github.com/genda/genda-api/internal/app"
	"github.com/genda/genda-api/internal/middlewares"
	"github.com/genda/genda-api/pkg/healthcheck"
)

// API constructs an http.Handler with all application routes defined.
func API(build string, shutdown chan os.Signal, log *log.Logger, postgresDB *sql.DB) http.Handler {

	healthHandler := healthcheck.NewHandler()
	basePermissions := []string{"genda-super-admin", "genda-staff"}

	a := app.New(shutdown, middlewares.Logger(log))

	// health
	a.Handle(http.MethodGet, "/health", healthHandler.GetHealthStatus)

	a = routes.UserRoutes(a, postgresDB, basePermissions)
	a = routes.StoreRoutes(a, postgresDB, basePermissions)
	return a
}
