package stores

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/genda/genda-api/internal/app"
	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

type handler struct {
	service    Service
	repository *StoreRepo
}

func NewHandler(postgresDB *sql.DB) *handler {
	storeRepository := NewStoreRepository(postgresDB)
	storeService := NewService(storeRepository)

	return &handler{
		service:    storeService,
		repository: storeRepository,
	}
}

// POST /stores
func (h *handler) CreateStore(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	var store Store
	if err := json.NewDecoder(r.Body).Decode(&store); err != nil {
		transformError(w, "Invalid request body", err.Error())
		return nil
	}

	validate := validator.New()
	if err := validate.Struct(store); err != nil {
		transformError(w, "Invalid request body", err.Error())
		return nil
	}

	res, err := h.service.CreateStore(store)
	if err != nil {
		transformError(w, "Failed to create store", err.Error())
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(res)
	return nil
}

// GET /stores
func (h *handler) GetStores(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	query := r.URL.Query()

	page, _ := strconv.Atoi(query.Get("page"))
	limit, _ := strconv.Atoi(query.Get("limit"))
	name := query.Get("name")
	storeType := query.Get("storeType")
	storeId := query.Get("storeId")

	stores, err := h.service.GetStores(page, limit, name, storeType, storeId)
	if err != nil {
		transformError(w, "Failed to get stores", err.Error())
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(stores)
	return nil
}

// GET /store/{id}
func (h *handler) GetStore(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	id := p.ByName("id")

	store, err := h.service.GetStore(id)
	if err != nil {
		transformError(w, "Failed to get store", err.Error())
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(store)
	return nil
}

// PUT /store/{id}
func (h *handler) UpdateStore(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	id := p.ByName("id")

	var store Store
	if err := json.NewDecoder(r.Body).Decode(&store); err != nil {
		transformError(w, "Invalid request body", err.Error())
		return nil
	}

	validate := validator.New()
	if err := validate.Struct(store); err != nil {
		transformError(w, "Invalid request body", err.Error())
		return nil
	}

	res, err := h.service.UpdateStore(id, store)
	if err != nil {
		transformError(w, "Failed to update store", err.Error())
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
	return nil
}

// DELETE /store/{id}
func (h *handler) DeleteStore(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	id := p.ByName("id")

	if err := h.service.DeleteStore(id); err != nil {
		transformError(w, "Failed to delete store", err.Error())
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	_ = json.NewEncoder(w).Encode("Store deleted")
	return nil
}

// POST /stores/plans
func (h *handler) CreateStorePlan(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	var plan StorePlan
	if err := json.NewDecoder(r.Body).Decode(&plan); err != nil {
		transformError(w, "Invalid request body", err.Error())
		return nil
	}

	validate := validator.New()
	if err := validate.Struct(plan); err != nil {
		transformError(w, "Invalid request body", err.Error())
		return nil
	}

	res, err := h.service.CreateStorePlan(plan)
	if err != nil {
		transformError(w, "Failed to create store plan", err.Error())
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(res)
	return nil
}

// GET /stores/plans?storeId={storeId}
func (h *handler) GetStorePlans(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	query := r.URL.Query()
	storeId := query.Get("storeId")

	plans, err := h.service.GetStorePlans(storeId)
	if err != nil {
		transformError(w, "Failed to get store plans", err.Error())
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(plans)
	return nil
}

// PUT /stores/plans/{id}
func (h *handler) UpdateStorePlan(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	id := p.ByName("id")

	var plan StorePlan
	if err := json.NewDecoder(r.Body).Decode(&plan); err != nil {
		transformError(w, "Invalid request body", err.Error())
		return nil
	}

	validate := validator.New()
	if err := validate.Struct(plan); err != nil {
		transformError(w, "Invalid request body", err.Error())
		return nil
	}

	res, err := h.service.UpdateStorePlan(id, plan)
	if err != nil {
		transformError(w, "Failed to update store plan", err.Error())
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
	return nil
}

// DELETE /stores/plans/{id}
func (h *handler) DeleteStorePlan(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	id := p.ByName("id")

	if err := h.service.DeleteStorePlan(id); err != nil {
		transformError(w, "Failed to delete store plan", err.Error())
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	_ = json.NewEncoder(w).Encode("Store plan deleted")
	return nil
}

// POST /stores/availability
func (h *handler) CreateStoreAvailability(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	var availability StoreAvailability
	if err := json.NewDecoder(r.Body).Decode(&availability); err != nil {
		transformError(w, "Invalid request body", err.Error())
		return nil
	}

	validate := validator.New()
	if err := validate.Struct(availability); err != nil {
		transformError(w, "Invalid request body", err.Error())
		return nil
	}

	res, err := h.service.CreateStoreAvailability(availability)
	if err != nil {
		transformError(w, "Failed to create store availability", err.Error())
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(res)
	return nil
}

// GET /stores/availability?storeId={storeId}
func (h *handler) GetStoreAvailability(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	query := r.URL.Query()
	storeId := query.Get("storeId")

	availability, err := h.service.GetStoreAvailability(storeId)
	if err != nil {
		transformError(w, "Failed to get store availability", err.Error())
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(availability)
	return nil
}

// PUT /stores/availability/{id}
func (h *handler) UpdateStoreAvailability(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	id := p.ByName("id")

	var availability StoreAvailability
	if err := json.NewDecoder(r.Body).Decode(&availability); err != nil {
		transformError(w, "Invalid request body", err.Error())
		return nil
	}

	validate := validator.New()
	if err := validate.Struct(availability); err != nil {
		transformError(w, "Invalid request body", err.Error())
		return nil
	}

	res, err := h.service.UpdateStoreAvailability(id, availability)
	if err != nil {
		transformError(w, "Failed to update store availability", err.Error())
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
	return nil
}

// DELETE /stores/availability/{id}
func (h *handler) DeleteStoreAvailability(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	id := p.ByName("id")

	if err := h.service.DeleteStoreAvailability(id); err != nil {
		transformError(w, "Failed to delete store availability", err.Error())
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	_ = json.NewEncoder(w).Encode("Store availability deleted")
	return nil
}

// POST /stores/ratings
func (h *handler) CreateStoreRating(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	var rating StoreRating
	if err := json.NewDecoder(r.Body).Decode(&rating); err != nil {
		transformError(w, "Invalid request body", err.Error())
		return nil
	}

	validate := validator.New()
	if err := validate.Struct(rating); err != nil {
		transformError(w, "Invalid request body", err.Error())
		return nil
	}

	res, err := h.service.CreateStoreRating(rating)
	if err != nil {
		transformError(w, "Failed to create store rating", err.Error())
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(res)
	return nil
}

// GET /stores/ratings?storeId={storeId}&page={page}&limit={limit}
func (h *handler) GetStoreRatings(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	query := r.URL.Query()

	storeId := query.Get("storeId")
	page, _ := strconv.Atoi(query.Get("page"))
	limit, _ := strconv.Atoi(query.Get("limit"))

	ratings, err := h.service.GetStoreRatings(storeId, page, limit)
	if err != nil {
		transformError(w, "Failed to get store ratings", err.Error())
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(ratings)
	return nil
}

// PUT /stores/ratings/{id}
func (h *handler) UpdateStoreRating(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	id := p.ByName("id")

	var rating StoreRating
	if err := json.NewDecoder(r.Body).Decode(&rating); err != nil {
		transformError(w, "Invalid request body", err.Error())
		return nil
	}

	validate := validator.New()
	if err := validate.Struct(rating); err != nil {
		transformError(w, "Invalid request body", err.Error())
		return nil
	}

	res, err := h.service.UpdateStoreRating(id, rating)
	if err != nil {
		transformError(w, "Failed to update store rating", err.Error())
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
	return nil
}

// DELETE /stores/ratings/{id}
func (h *handler) DeleteStoreRating(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	id := p.ByName("id")

	if err := h.service.DeleteStoreRating(id); err != nil {
		transformError(w, "Failed to delete store rating", err.Error())
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	_ = json.NewEncoder(w).Encode("Store rating deleted")
	return nil
}

// POST /stores/appointments
func (h *handler) CreateStoreAppointment(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	var appointment StoreAppointment
	if err := json.NewDecoder(r.Body).Decode(&appointment); err != nil {
		transformError(w, "Invalid request body", err.Error())
		return nil
	}

	validate := validator.New()
	if err := validate.Struct(appointment); err != nil {
		transformError(w, "Invalid request body", err.Error())
		return nil
	}

	res, err := h.service.CreateStoreAppointment(appointment)
	if err != nil {
		transformError(w, "Failed to create store appointment", err.Error())
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(res)
	return nil
}

// GET /stores/appointments?storeId={storeId}&page={page}&limit={limit}
func (h *handler) GetStoreAppointments(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	query := r.URL.Query()
	storeId := query.Get("storeId")
	page, _ := strconv.Atoi(query.Get("page"))
	limit, _ := strconv.Atoi(query.Get("limit"))

	appointments, err := h.service.GetStoreAppointments(storeId, page, limit)
	if err != nil {
		transformError(w, "Failed to get store appointments", err.Error())
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(appointments)
	return nil
}

// PUT /stores/appointments/{id}
func (h *handler) UpdateStoreAppointment(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	id := p.ByName("id")

	var appointment StoreAppointment
	if err := json.NewDecoder(r.Body).Decode(&appointment); err != nil {
		transformError(w, "Invalid request body", err.Error())
		return nil
	}

	validate := validator.New()
	if err := validate.Struct(appointment); err != nil {
		transformError(w, "Invalid request body", err.Error())
		return nil
	}

	res, err := h.service.UpdateStoreAppointment(id, appointment)
	if err != nil {
		transformError(w, "Failed to update store appointment", err.Error())
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
	return nil
}

// DELETE /stores/appointments/{id}
func (h *handler) DeleteStoreAppointment(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	id := p.ByName("id")

	if err := h.service.DeleteStoreAppointment(id); err != nil {
		transformError(w, "Failed to delete store appointment", err.Error())
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	_ = json.NewEncoder(w).Encode("Store appointment deleted")
	return nil
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
