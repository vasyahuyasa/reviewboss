package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/vasyahuyasa/reviewboss/internal/core/review"
)

type ReviewHandlers struct {
	reviewService *review.Service
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

	mr, err := h.reviewService.RegisterMergeRequest(id, url, review.SkillGolang)
	switch {
	case err == review.ErrNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
		return

	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	default:
		log.Println("merger request", id, "registered")
	}

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

	reviwers, err := h.reviewService.ReviwersBySkill(skill)

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
