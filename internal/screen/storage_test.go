package screen

import (
	"sync"
	"testing"
)

func resetStore() {
	mu.Lock()
	defer mu.Unlock()
	screenStore = map[string]ScreenResponse{
		"home": newHomeScreen(),
	}
}

func TestGetScreenInitialHome(t *testing.T) {
	resetStore()
	resp, ok := getScreen("home")
	if !ok {
		t.Fatal("expected home screen to exist")
	}
	if resp.ScreenName != "home" {
		t.Errorf("expected screen_name home, got %s", resp.ScreenName)
	}
	if !resp.PullToRefresh {
		t.Error("expected pull_to_refresh true")
	}
	if len(resp.Components) != 2 {
		t.Errorf("expected 2 components, got %d", len(resp.Components))
	}
}

func TestStorageGetScreenNotFound(t *testing.T) {
	resetStore()
	_, ok := getScreen("nonexistent")
	if ok {
		t.Error("expected nonexistent screen to not be found")
	}
}

func TestSetAndGetScreen(t *testing.T) {
	resetStore()
	s := ScreenResponse{ScreenName: "test", PullToRefresh: false, Components: nil}
	setScreen("test", s)
	resp, ok := getScreen("test")
	if !ok {
		t.Fatal("expected test screen to exist")
	}
	if resp.ScreenName != "test" {
		t.Errorf("expected screen_name test, got %s", resp.ScreenName)
	}
}

func TestStorageDeleteScreen(t *testing.T) {
	resetStore()
	setScreen("test", ScreenResponse{ScreenName: "test"})
	if !deleteScreen("test") {
		t.Error("expected delete to return true")
	}
	if _, ok := getScreen("test"); ok {
		t.Error("expected test screen to be deleted")
	}
}

func TestStorageDeleteScreenNotFound(t *testing.T) {
	resetStore()
	if deleteScreen("nonexistent") {
		t.Error("expected delete of nonexistent to return false")
	}
}

func TestConcurrentAccess(t *testing.T) {
	resetStore()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			id := "concurrent"
			setScreen(id, ScreenResponse{ScreenName: id})
			getScreen(id)
			deleteScreen(id)
		}(i)
	}
	wg.Wait()
}
