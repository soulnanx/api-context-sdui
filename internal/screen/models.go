package screen

import "encoding/json"

type Component struct {
	Type string          `json:"type"`
	ID   string          `json:"id"`
	Data json.RawMessage `json:"data"`
}

type ScreenResponse struct {
	ScreenName    string      `json:"screen_name"`
	PullToRefresh bool        `json:"pull_to_refresh"`
	Components    []Component `json:"components"`
}
