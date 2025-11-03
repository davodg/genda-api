package subscriptions

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type handler struct {
	service    Service
	repository *SubscriptionRepo
}

func NewHandler(postgresDB *sql.DB) *handler {
	subscriptionRepository := NewSubscriptionRepository(postgresDB)
	subscriptionService := NewService(subscriptionRepository)

	return &handler{
		service:    subscriptionService,
		repository: subscriptionRepository,
	}
}

// POST /subscriptions
func (h *handler) CreateSubscription(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	var subscription Subscription
	if err := json.NewDecoder(r.Body).Decode(&subscription); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return err
	}

	createdSubscription, err := h.service.CreateSubscription(subscription)
	if err != nil {
		http.Error(w, "Failed to create subscription", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(createdSubscription)
}

// GET /subscriptions?user_id={userId}&page={page}&limit={limit}
func (h *handler) GetSubscriptions(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	userId := r.URL.Query().Get("user_id")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	subscriptionsResponse, err := h.service.GetSubscriptions(userId, page, limit)
	if err != nil {
		http.Error(w, "Failed to get subscriptions", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(subscriptionsResponse)
}

// PUT /subscriptions/{id}
func (h *handler) UpdateSubscription(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	id := p.ByName("id")
	var subscription Subscription
	if err := json.NewDecoder(r.Body).Decode(&subscription); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return err
	}

	updatedSubscription, err := h.service.UpdateSubscription(id, subscription)
	if err != nil {
		http.Error(w, "Failed to update subscription", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(updatedSubscription)
}

// DELETE /subscriptions/{id}
func (h *handler) DeleteSubscription(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	id := p.ByName("id")

	if err := h.service.DeleteSubscription(id); err != nil {
		http.Error(w, "Failed to delete subscription", http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
