package healthcheck

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type handler struct{}

func NewHandler() *handler {
	return &handler{}
}

// GET /health
func (h *handler) GetHealthStatus(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {

	data := map[string]string{
		"healthy": "OK",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
	return nil
}
