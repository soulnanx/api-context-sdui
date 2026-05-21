package component

type HeroBanner struct {
	Title    string `json:"title"`
	ImageURL string `json:"image_url"`
	ActionID string `json:"action_id"`
}

type ActionButton struct {
	Text     string `json:"text"`
	ActionID string `json:"action_id"`
}
