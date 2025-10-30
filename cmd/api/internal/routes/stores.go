package routes

import (
	"database/sql"
	"net/http"

	"github.com/genda/genda-api/internal/app"
	"github.com/genda/genda-api/internal/middlewares"
	"github.com/genda/genda-api/pkg/stores"
)

func StoreRoutes(a *app.App, postgresDB *sql.DB, basePermissions []string) *app.App {

	organizationHandler := stores.NewHandler(postgresDB)

	a.Handle(http.MethodPost, "/api/v1/stores", organizationHandler.CreateStore, middlewares.Authenticate(append(basePermissions, []string{"genda-admin", "genda-associated"}...)))
	a.Handle(http.MethodGet, "/api/v1/stores", organizationHandler.GetStores, middlewares.Authenticate(append(basePermissions, []string{"genda-admin", "genda-associated", "genda-analyst", "genda-consultant"}...)))
	a.Handle(http.MethodGet, "/api/v1/stores/:id", organizationHandler.GetStore, middlewares.Authenticate(append(basePermissions, []string{"genda-admin", "genda-associated", "genda-analyst", "genda-consultant"}...)))
	a.Handle(http.MethodPut, "/api/v1/stores/:id", organizationHandler.UpdateStore, middlewares.Authenticate(append(basePermissions, []string{"genda-admin", "genda-associated"}...)))
	a.Handle(http.MethodDelete, "/api/v1/stores/:id", organizationHandler.DeleteStore, middlewares.Authenticate(append(basePermissions, []string{"genda-admin", "genda-associated"}...)))

	return a
}
