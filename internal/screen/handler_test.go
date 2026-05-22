package screen

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestServer() *httptest.Server {
	mux := http.NewServeMux()
	screenHandler := &ScreenHandler{}
	mux.HandleFunc("GET /api/screen/{screenId}", screenHandler.GetScreen)
	mux.HandleFunc("POST /api/admin/screen/{screenId}", screenHandler.CreateScreen)
	mux.HandleFunc("GET /api/admin/screen/{screenId}", screenHandler.GetScreenConfig)
	mux.HandleFunc("PUT /api/admin/screen/{screenId}", screenHandler.UpdateScreen)
	mux.HandleFunc("DELETE /api/admin/screen/{screenId}", screenHandler.DeleteScreen)
	return httptest.NewServer(mux)
}

func TestGetScreenHome(t *testing.T) {
	resetStore()
	ts := setupTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/screen/home")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	var sr ScreenResponse
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if sr.ScreenName != "home" {
		t.Errorf("expected screen_name home, got %s", sr.ScreenName)
	}
}

func TestGetScreenNotFound(t *testing.T) {
	resetStore()
	ts := setupTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/screen/nonexistent")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestCreateScreen(t *testing.T) {
	resetStore()
	ts := setupTestServer()
	defer ts.Close()

	body := []byte(`{"pull_to_refresh":true,"components":[{"type":"hero_banner","id":"b1","data":{"title":"Test"}}]}`)
	resp, err := http.Post(ts.URL+"/api/admin/screen/test", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}

	// Verify it exists via GET /api/screen/test
	resp2, err := http.Get(ts.URL + "/api/screen/test")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Errorf("expected 200 after create, got %d", resp2.StatusCode)
	}
}

func TestCreateScreenDuplicate(t *testing.T) {
	resetStore()
	ts := setupTestServer()
	defer ts.Close()

	body := []byte(`{"pull_to_refresh":true,"components":[]}`)
	resp, err := http.Post(ts.URL+"/api/admin/screen/home", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusConflict {
		t.Errorf("expected 409, got %d", resp.StatusCode)
	}
}

func TestGetScreenConfig(t *testing.T) {
	resetStore()
	ts := setupTestServer()
	defer ts.Close()

	// Create first
	body := []byte(`{"pull_to_refresh":false,"components":[]}`)
	resp, err := http.Post(ts.URL+"/api/admin/screen/test", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	resp.Body.Close()

	// Get config
	resp2, err := http.Get(ts.URL + "/api/admin/screen/test")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp2.StatusCode)
	}
}

func TestGetScreenConfigNotFound(t *testing.T) {
	resetStore()
	ts := setupTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/admin/screen/nonexistent")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestUpdateScreen(t *testing.T) {
	resetStore()
	ts := setupTestServer()
	defer ts.Close()

	// Create first
	body := []byte(`{"pull_to_refresh":true,"components":[]}`)
	resp, err := http.Post(ts.URL+"/api/admin/screen/test", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	resp.Body.Close()

	// Update
	body2 := []byte(`{"pull_to_refresh":false,"components":[{"type":"action_button","id":"btn1","data":{"text":"Test"}}]}`)
	req, _ := http.NewRequest(http.MethodPut, ts.URL+"/api/admin/screen/test", bytes.NewBuffer(body2))
	req.Header.Set("Content-Type", "application/json")
	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp2.StatusCode)
	}

	// Verify update
	resp3, err := http.Get(ts.URL + "/api/screen/test")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp3.Body.Close()

	if resp3.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp3.StatusCode)
	}
	var sr ScreenResponse
	json.NewDecoder(resp3.Body).Decode(&sr)
	if sr.PullToRefresh {
		t.Error("expected pull_to_refresh false after update")
	}
}

func TestUpdateScreenNotFound(t *testing.T) {
	resetStore()
	ts := setupTestServer()
	defer ts.Close()

	body := []byte(`{"pull_to_refresh":true,"components":[]}`)
	req, _ := http.NewRequest(http.MethodPut, ts.URL+"/api/admin/screen/nonexistent", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestDeleteScreen(t *testing.T) {
	resetStore()
	ts := setupTestServer()
	defer ts.Close()

	// Create first
	body := []byte(`{"pull_to_refresh":true,"components":[]}`)
	resp, err := http.Post(ts.URL+"/api/admin/screen/test", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	resp.Body.Close()

	// Delete
	req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/api/admin/screen/test", nil)
	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusNoContent {
		t.Errorf("expected 204, got %d", resp2.StatusCode)
	}

	// Verify deleted
	resp3, err := http.Get(ts.URL + "/api/screen/test")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp3.Body.Close()

	if resp3.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 after delete, got %d", resp3.StatusCode)
	}
}

func TestDeleteScreenNotFound(t *testing.T) {
	resetStore()
	ts := setupTestServer()
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/api/admin/screen/nonexistent", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

// Helper to read body (unused but kept for reference)
var _ = io.ReadAll
