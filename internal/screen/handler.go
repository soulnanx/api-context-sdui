package screen

import (
	"encoding/json"
	"log"
	"net/http"
)

type ScreenHandler struct{}

func (h *ScreenHandler) GetScreen(w http.ResponseWriter, r *http.Request) {
	screenId := r.PathValue("screenId")
	if screenId == "" {
		http.Error(w, "missing screenId", http.StatusBadRequest)
		return
	}

	resp, ok := getScreen(screenId)
	if !ok {
		http.Error(w, "screen not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (h *ScreenHandler) CreateScreen(w http.ResponseWriter, r *http.Request) {
	screenId := r.PathValue("screenId")
	if screenId == "" {
		http.Error(w, "missing screenId", http.StatusBadRequest)
		return
	}

	if _, ok := getScreen(screenId); ok {
		http.Error(w, "screen already exists", http.StatusConflict)
		return
	}

	var req ScreenResponse
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	req.ScreenName = screenId

	setScreen(screenId, req)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(req); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (h *ScreenHandler) GetScreenConfig(w http.ResponseWriter, r *http.Request) {
	screenId := r.PathValue("screenId")
	if screenId == "" {
		http.Error(w, "missing screenId", http.StatusBadRequest)
		return
	}

	resp, ok := getScreen(screenId)
	if !ok {
		http.Error(w, "screen not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (h *ScreenHandler) UpdateScreen(w http.ResponseWriter, r *http.Request) {
	screenId := r.PathValue("screenId")
	if screenId == "" {
		http.Error(w, "missing screenId", http.StatusBadRequest)
		return
	}

	if _, ok := getScreen(screenId); !ok {
		http.Error(w, "screen not found", http.StatusNotFound)
		return
	}

	var req ScreenResponse
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	req.ScreenName = screenId

	setScreen(screenId, req)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(req); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (h *ScreenHandler) DeleteScreen(w http.ResponseWriter, r *http.Request) {
	screenId := r.PathValue("screenId")
	if screenId == "" {
		http.Error(w, "missing screenId", http.StatusBadRequest)
		return
	}

	if !deleteScreen(screenId) {
		http.Error(w, "screen not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
