package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"nexturn_final/internal/logger"
	"nexturn_final/internal/service"
)

type Handler struct {
	svc  *service.URLService
	logg *logger.Logger
}

func NewHandler(svc *service.URLService, logg *logger.Logger) *Handler {
	return &Handler{svc: svc, logg: logg}
}

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL     string    `json:"url"`
		Expires time.Time `json:"expires,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if req.Expires.IsZero() {
		req.Expires = time.Now().Add(24 * time.Hour)
	}
	ip := r.RemoteAddr
	code, err := h.svc.Shorten(r.Context(), req.URL, req.Expires, ip)
	if err != nil {
		http.Error(w, "could not shorten", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"code": code})
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	ip := r.RemoteAddr
	url, err := h.svc.Resolve(r.Context(), code, ip)
	if err != nil {
		http.Error(w, "not found or expired", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	hits, createdAt, expiresAt, err := h.svc.Stats(r.Context(), code)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"hits":      hits,
		"createdAt": createdAt,
		"expiresAt": expiresAt,
	})
}
