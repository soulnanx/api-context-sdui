package screen

import (
	"encoding/json"
	"log"
	"net/http"

	"api-context-sdui/internal/component"
)

type ScreenHandler struct{}

func (h *ScreenHandler) GetHomeScreen(w http.ResponseWriter, r *http.Request) {
	hero := component.HeroBanner{
		Title:    "Descubra o SDUI no Android",
		ImageURL: "https://images.unsplash.com/photo-1607604276583-eef5d076aa5f",
		ActionID: "NAVIGATE_TO_DETAILS",
	}
	heroData, err := json.Marshal(hero)
	if err != nil {
		log.Printf("failed to marshal hero banner: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	btn := component.ActionButton{
		Text:     "Clique para Recarregar",
		ActionID: "RELOAD_SCREEN",
	}
	btnData, err := json.Marshal(btn)
	if err != nil {
		log.Printf("failed to marshal action button: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := ScreenResponse{
		ScreenName:    "home_page",
		PullToRefresh: true,
		Components: []Component{
			{Type: "hero_banner", ID: "b1", Data: heroData},
			{Type: "action_button", ID: "btn1", Data: btnData},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
