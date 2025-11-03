package routes

import (
	"database/sql"
	"net/http"

	"github.com/genda/genda-api/internal/app"
	"github.com/genda/genda-api/internal/middlewares"
	"github.com/genda/genda-api/pkg/users"
)

func UserRoutes(a *app.App, postgresDB *sql.DB, basePermissions []string) *app.App {

	userHandler := users.NewHandler(postgresDB)

	a.Handle(http.MethodPost, "/api/v1/users", userHandler.CreateUser, middlewares.Authenticate(append(basePermissions, []string{"genda-owner", "genda-admin", "genda-customer"}...)))
	a.Handle(http.MethodGet, "/api/v1/users", userHandler.GetUsers, middlewares.Authenticate(append(basePermissions, []string{"genda-owner", "genda-admin", "genda-customer"}...)))
	a.Handle(http.MethodGet, "/api/v1/users/:id", userHandler.GetUser, middlewares.Authenticate(append(basePermissions, []string{"genda-owner", "genda-admin", "genda-customer"}...)))
	a.Handle(http.MethodPut, "/api/v1/users/:id", userHandler.UpdateUser, middlewares.Authenticate(append(basePermissions, []string{"genda-owner", "genda-admin", "genda-customer"}...)))
	a.Handle(http.MethodDelete, "/api/v1/users/:id", userHandler.DeleteUser, middlewares.Authenticate(append(basePermissions, []string{"genda-admin"}...)))

	a.Handle(http.MethodPost, "/api/v1/auth/login", userHandler.AuthUser)
	a.Handle(http.MethodPost, "/api/v1/auth/logout", userHandler.LogoutUser)

	return a
}
