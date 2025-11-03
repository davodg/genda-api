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

	a.Handle(http.MethodPost, "/api/v1/stores", organizationHandler.CreateStore, middlewares.Authenticate(append(basePermissions, []string{"genda-admin", "genda-owner"}...)))
	a.Handle(http.MethodGet, "/api/v1/stores", organizationHandler.GetStores, middlewares.Authenticate(append(basePermissions, []string{"genda-admin", "genda-owner", "genda-customer"}...)))
	a.Handle(http.MethodGet, "/api/v1/stores/:id", organizationHandler.GetStore, middlewares.Authenticate(append(basePermissions, []string{"genda-admin", "genda-owner", "genda-customer"}...)))
	a.Handle(http.MethodPut, "/api/v1/stores/:id", organizationHandler.UpdateStore, middlewares.Authenticate(append(basePermissions, []string{"genda-admin", "genda-owner"}...)))
	a.Handle(http.MethodDelete, "/api/v1/stores/:id", organizationHandler.DeleteStore, middlewares.Authenticate(append(basePermissions, []string{"genda-admin", "genda-owner"}...)))

	a.Handler(http.MethodPost, "/api/v1/stores/:id/plans", organizationHandler.CreateStorePlan, middlewares.Authenticate(append(basePermissions, []string{"genda-admin", "genda-owner"}...)))
	a.Handler(http.MethodGet, "/api/v1/stores/:id/plans", organizationHandler.GetStorePlans, middlewares.Authenticate(append(basePermissions, []string{"genda-admin", "genda-owner", "genda-customer"}...)))
	a.Handler(http.MethodPut, "/api/v1/stores/:id/plans/:planId", organizationHandler.UpdateStorePlan, middlewares.Authenticate(append(basePermissions, []string{"genda-admin", "genda-owner"}...)))
	a.Handler(http.MethodDelete, "/api/v1/stores/:id/plans/:planId", organizationHandler.DeleteStorePlan, middlewares.Authenticate(append(basePermissions, []string{"genda-admin", "genda-owner"}...)))

	a.Handler(http.MethodPost, "/api/v1/stores/:id/availability", organizationHandler.CreateStoreAvailability, middlewares.Authenticate(append(basePermissions, []string{"genda-admin", "genda-owner"}...)))
	a.Handler(http.MethodGet, "/api/v1/stores/:id/availability", organizationHandler.GetStoreAvailability, middlewares.Authenticate(append(basePermissions, []string{"genda-admin", "genda-owner", "genda-customer"}...)))
	a.Handler(http.MethodPut, "/api/v1/stores/:id/availability/:availabilityId", organizationHandler.UpdateStoreAvailability, middlewares.Authenticate(append(basePermissions, []string{"genda-admin", "genda-owner"}...)))
	a.Handler(http.MethodDelete, "/api/v1/stores/:id/availability/:availabilityId", organizationHandler.DeleteStoreAvailability, middlewares.Authenticate(append(basePermissions, []string{"genda-admin", "genda-owner"}...)))

	a.Handler(http.MethodPost, "/api/v1/stores/:id/ratings", organizationHandler.CreateStoreRating, middlewares.Authenticate(append(basePermissions, []string{"genda-admin", "genda-owner", "genda-customer"}...)))
	a.Handler(http.MethodGet, "/api/v1/stores/:id/ratings", organizationHandler.GetStoreRatings, middlewares.Authenticate(append(basePermissions, []string{"genda-admin", "genda-owner", "genda-customer"}...)))
	a.Handler(http.MethodPut, "/api/v1/stores/:id/ratings/:ratingId", organizationHandler.UpdateStoreRating, middlewares.Authenticate(append(basePermissions, []string{"genda-admin", "genda-owner", "genda-customer"}...)))
	a.Handler(http.MethodDelete, "/api/v1/stores/:id/ratings/:ratingId", organizationHandler.DeleteStoreRating, middlewares.Authenticate(append(basePermissions, []string{"genda-admin", "genda-owner", "genda-customer"}...)))

	a.Handler(http.MethodPost, "/api/v1/stores/:id/appointments", organizationHandler.CreateStoreAppointment, middlewares.Authenticate(append(basePermissions, []string{"genda-admin", "genda-owner", "genda-customer"}...)))
	a.Handler(http.MethodGet, "/api/v1/stores/:id/appointments", organizationHandler.GetStoreAppointments, middlewares.Authenticate(append(basePermissions, []string{"genda-admin", "genda-owner", "genda-customer"}...)))
	a.Handler(http.MethodPut, "/api/v1/stores/:id/appointments/:appointmentId", organizationHandler.UpdateStoreAppointment, middlewares.Authenticate(append(basePermissions, []string{"genda-admin", "genda-owner", "genda-customer"}...)))
	a.Handler(http.MethodDelete, "/api/v1/stores/:id/appointments/:appointmentId", organizationHandler.DeleteStoreAppointment, middlewares.Authenticate(append(basePermissions, []string{"genda-admin", "genda-owner", "genda-customer"}...)))

	return a
}
