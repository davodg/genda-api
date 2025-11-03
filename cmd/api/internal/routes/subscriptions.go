package routes

import (
	"database/sql"
	"net/http"

	"github.com/genda/genda-api/internal/app"
	"github.com/genda/genda-api/internal/middlewares"
	"github.com/genda/genda-api/pkg/subscriptions"
)

func SubscriptionRoutes(a *app.App, postgresDB *sql.DB, basePermissions []string) *app.App {

	subscriptionHandler := subscriptions.NewHandler(postgresDB)

	a.Handle(http.MethodPost, "/api/v1/subscriptions", subscriptionHandler.CreateSubscription, middlewares.Authenticate(append(basePermissions, []string{"genda-admin", "genda-owner", "genda-customer"}...)))
	a.Handle(http.MethodGet, "/api/v1/subscriptions", subscriptionHandler.GetSubscriptions, middlewares.Authenticate(append(basePermissions, []string{"genda-admin", "genda-owner", "genda-customer"}...)))
	a.Handle(http.MethodPut, "/api/v1/subscriptions/:id", subscriptionHandler.UpdateSubscription, middlewares.Authenticate(append(basePermissions, []string{"genda-admin", "genda-owner"}...)))
	a.Handle(http.MethodDelete, "/api/v1/subscriptions/:id", subscriptionHandler.DeleteSubscription, middlewares.Authenticate(append(basePermissions, []string{"genda-admin"}...)))

	return a
}
