package inbound

import (
	"encoding/json"
	"net/http"

	"github.com/vasyahuyasa/reviewboss/internal/domain"
	"github.com/vasyahuyasa/reviewboss/internal/ports"
)

type HTTPSource struct {
	Store ports.MergeRequestSource
}

func (h *HTTPSource) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var mr domain.MergeRequest
	if err := json.NewDecoder(r.Body).Decode(&mr); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	h.Store.Save(&mr)
	w.WriteHeader(http.StatusAccepted)
}
