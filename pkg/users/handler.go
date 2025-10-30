package users

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/genda/genda-api/internal/app"
	"github.com/golang-jwt/jwt/v4"

	"github.com/fesa/genda-api/pkg/config"
	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

type handler struct {
	service    Service
	repository *UserRepo
}

func NewHandler(postgresDB *sql.DB) *handler {

	userRepository := NewUserRepository(postgresDB)
	userService := NewService(userRepository)

	return &handler{
		service:    userService,
		repository: userRepository,
	}
}

// POST /users
func (h *handler) CreateUser(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		transformError(w, "Invalid request body", err.Error())
		return nil
	}

	validate := validator.New()
	err = validate.Struct(user)
	if err != nil {
		transformError(w, "Invalid request body", err.Error())
		return nil
	}

	res, err := h.service.CreateUser(user)
	if err != nil {
		transformError(w, "Failed to create user", err.Error())
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
	return nil
}

// GET /users
func (h *handler) GetUsers(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	queryParams := r.URL.Query()

	page, _ := strconv.Atoi(queryParams.Get("page"))
	limit, _ := strconv.Atoi(queryParams.Get("limit"))
	name := queryParams.Get("name")
	email := queryParams.Get("email")
	userId := queryParams.Get("userId")

	users, err := h.service.GetUsers(page, limit, name, email, userId)
	if err != nil {
		transformError(w, "Failed to get users", err.Error())
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
	return nil
}

// GET /user/{id}
func (h *handler) GetUser(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	id := p.ByName("id")

	user, err := h.service.GetUser(id)
	if err != nil {
		transformError(w, "Failed to get user", err.Error())
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
	return nil
}

// PUT /user/{id}
func (h *handler) UpdateUser(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	id := p.ByName("id")

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		transformError(w, "Invalid request body", err.Error())
		return nil
	}

	validate := validator.New()
	err = validate.Struct(user)
	if err != nil {
		transformError(w, "Invalid request body", err.Error())
		return nil
	}

	res, err := h.service.UpdateUser(id, user)
	if err != nil {
		transformError(w, "Failed to update user", err.Error())
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
	return nil
}

// DELETE /user/{id}
func (h *handler) DeleteUser(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	id := p.ByName("id")

	err := h.service.DeleteUser(id)
	if err != nil {
		transformError(w, "Failed to delete user", err.Error())
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(204)
	json.NewEncoder(w).Encode("User deleted")
	return nil
}

// POST /users/login
func (i *handler) AuthUser(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	conf := config.New()
	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return nil
	}

	data := url.Values{}
	data.Set("client_id", conf.KeycloakClient)
	data.Set("client_secret", conf.ClientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("code", loginReq.Code)
	data.Set("redirect_uri", loginReq.RedirectUrl)
	data.Set("scope", "openid broker")

	// data.Set("username", loginReq.UserName)
	// data.Set("password", loginReq.Password)

	host := conf.KeycloakHost + "/protocol/openid-connect/token"

	req, err := http.NewRequest("POST", host, strings.NewReader(data.Encode()))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return nil
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to connect to Keycloak", http.StatusInternalServerError)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("Error, ", string(body))
		http.Error(w, fmt.Sprintf("Keycloak error: %s", string(body)), resp.StatusCode)
		return nil
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		http.Error(w, "Failed to parse token response", http.StatusInternalServerError)
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenResp)
	return nil
}

func (i *handler) LogoutUser(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	conf := config.New()
	var logoutReq LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&logoutReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return nil
	}

	jwtToken, err := i.decodeJwt(logoutReq.Token)
	if err != nil {
		log.Println("error: ", err)
		http.Error(w, "Error decoding token", http.StatusBadRequest)
		return fmt.Errorf("error decoding token: %v", err)
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("failed to convert claims")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return fmt.Errorf("sub claim not found or not a string")
	}

	accessToken, err := i.getAccessToken()
	if err != nil {
		http.Error(w, "Error getting keycloak access token", http.StatusBadRequest)
		return nil
	}

	host := conf.KeycloakHostWithoutRealm + "/admin/realms/" + conf.KeycloakRealm + "/users/" + sub + "/logout"

	req, err := http.NewRequest("POST", host, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return nil
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to connect to Keycloak", http.StatusInternalServerError)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return nil
	}

	if resp.StatusCode != http.StatusNoContent {
		log.Println("Error, ", string(body))
		http.Error(w, fmt.Sprintf("Keycloak error: %s", string(body)), resp.StatusCode)
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode("Logged out")
	return nil
}

func (i *handler) getAccessToken() (string, error) {
	conf := config.New()

	data := url.Values{}
	data.Set("client_id", conf.KeycloakClient)
	data.Set("client_secret", conf.ClientSecret)
	data.Set("grant_type", "client_credentials")

	host := conf.KeycloakHost + "/protocol/openid-connect/token"

	req, err := http.NewRequest("POST", host, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %v", err)
	}

	return result.AccessToken, nil
}

func (i *handler) decodeJwt(tokenString string) (*jwt.Token, error) {

	parser := new(jwt.Parser)
	token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// transform error for response api
func transformError(w http.ResponseWriter, m string, e string) {
	var data = app.ValidateError{
		Message: m,
		Error:   e,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(data)
}
