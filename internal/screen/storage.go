package screen

import (
	"encoding/json"
	"log"
	"sync"

	"api-context-sdui/internal/component"
)

var (
	mu          sync.RWMutex
	screenStore = make(map[string]ScreenResponse)
)

func init() {
	screenStore["home"] = newHomeScreen()
}

func newHomeScreen() ScreenResponse {
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

func getScreen(id string) (ScreenResponse, bool) {
	mu.RLock()
	defer mu.RUnlock()
	s, ok := screenStore[id]
	return s, ok
}

func setScreen(id string, resp ScreenResponse) {
	mu.Lock()
	defer mu.Unlock()
	screenStore[id] = resp
}

func deleteScreen(id string) bool {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := screenStore[id]; !ok {
		return false
	}
	delete(screenStore, id)
	return true
}
