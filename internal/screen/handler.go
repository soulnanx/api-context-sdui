package screen

import (
	"encoding/json"
	"log"
	"net/http"

	"api-context-sdui/internal/component"
)

type ScreenHandler struct{}

func (h *ScreenHandler) GetScreen(w http.ResponseWriter, r *http.Request) {
	screenId := r.PathValue("screenId")
	if screenId == "" {
		http.Error(w, "missing screenId", http.StatusBadRequest)
		return
	}

	resp, ok := h.buildScreen(screenId)
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

func (h *ScreenHandler) buildScreen(screenId string) (ScreenResponse, bool) {
	switch screenId {
	case "home":
		return h.buildHomeScreen(), true
	default:
		return ScreenResponse{}, false
	}
}

func (h *ScreenHandler) buildHomeScreen() ScreenResponse {
	hero := component.HeroBanner{
		Title:    "Descubra o SDUI no Android",
		ImageURL: "https://images.unsplash.com/photo-1607604276583-eef5d076aa5f",
		ActionID: "NAVIGATE_TO_DETAILS",
	}
	heroData, err := json.Marshal(hero)
	if err != nil {
		log.Printf("failed to marshal hero banner: %v", err)
		return ScreenResponse{}
	}

	btn := component.ActionButton{
		Text:     "Clique para Recarregar",
		ActionID: "RELOAD_SCREEN",
	}
	btnData, err := json.Marshal(btn)
	if err != nil {
		log.Printf("failed to marshal action button: %v", err)
		return ScreenResponse{}
	}

	return ScreenResponse{
		ScreenName:    "home",
		PullToRefresh: true,
		Components: []Component{
			{Type: "hero_banner", ID: "b1", Data: heroData},
			{Type: "action_button", ID: "btn1", Data: btnData},
		},
	}
}
