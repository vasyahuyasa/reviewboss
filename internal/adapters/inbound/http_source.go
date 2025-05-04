package inbound

import (
	"encoding/json"
	"net/http"

	"github.com/vasyahuyasa/reviewboss/internal/domain"
)

type HTTPSource struct {
	boss domain.Boss
}

func (h *HTTPSource) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var mr domain.MergeRequest
	if err := json.NewDecoder(r.Body).Decode(&mr); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	h.boss.AddMR(mr domain.MergeRequest)
	w.WriteHeader(http.StatusAccepted)
}
