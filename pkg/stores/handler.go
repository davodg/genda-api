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
