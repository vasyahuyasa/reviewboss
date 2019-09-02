package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/vasyahuyasa/reviewboss/internal/core/review"
)

type ReviewHandlers struct {
	Brain *review.Brain
}

func (h *ReviewHandlers) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `
	<html>
		<head></head>
		<body>
			<ul>
				<li><a href="/register">/register</a></li>
				<li><a href="/reviwers">/reviwers</a></li>
			</ul>
		</body>
	</html>
	`)
}

func (h *ReviewHandlers) Register(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "url required", http.StatusBadRequest)
		return
	}

	h.Brain.Lock()
	defer h.Brain.Unlock()

	skill := review.SkillGolang

	mr := review.NewMergeRequest(id, url, skill)

	err := h.Brain.RegisterMergeRequest(mr)
	if err != nil {
		log.Printf("can not register merge request: %v", err)
		http.Error(w, fmt.Sprintf("can not register merge request: %v", err), http.StatusInternalServerError)
		return
	}

	reviwers, err := h.Brain.SelectReviwers(skill)
	if err != nil {
		log.Printf("can not get list of reviwers: %v", err)
		http.Error(w, fmt.Sprintf("can not can not get list of reviwers: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("merger request", id, "registered")

	enc := json.NewEncoder(w)
	err = enc.Encode(mr)
	if err != nil {
		log.Printf("can not marshal mergerequest: %v", err)
		http.Error(w, fmt.Sprintf("can not marshal mergerqeuest: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *ReviewHandlers) ListReviwers(w http.ResponseWriter, r *http.Request) {
	skill := review.SkillGolang

	enc := json.NewEncoder(w)

	reviwers, err := h.Brain.SelectReviwers(skill)

	if err != nil {
		log.Printf("can not get list of reviwers: %v", err)
		http.Error(w, "can not get list of reviwers: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = enc.Encode(reviwers)
	if err != nil {
		log.Printf("can not encode list of reviwers: %v", err)
	}
}

func (h *ReviewHandlers) Accept(w http.ResponseWriter, r *http.Request) {

}

func (h *ReviewHandlers) Decline(w http.ResponseWriter, r *http.Request) {

}