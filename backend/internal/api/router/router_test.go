package router

import (
	"backend/internal/config"
	"backend/internal/utils"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterInitialization(t *testing.T) {
	cfg := &config.Config{
		Port:      "8080",
		GinMode:   "test",
		LogLevel:  "info",
		LogFormat: "text",
	}

	logger := slog.Default()
	geoManager := &utils.GeoJsonManager{}

	r := Setup(cfg, geoManager, logger)

	if r == nil {
		t.Fatal("router should not be nil")
	}
}

func TestHealthEndpoint(t *testing.T) {
	cfg := &config.Config{
		Port:      "8080",
		GinMode:   "test",
		LogLevel:  "info",
		LogFormat: "text",
	}

	logger := slog.Default()
	geoManager := &utils.GeoJsonManager{}

	r := Setup(cfg, geoManager, logger)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response["status"] != "up" {
		t.Errorf("expected status 'up', got %v", response["status"])
	}

	if _, ok := response["build_time"]; !ok {
		t.Error("response should contain build_time")
	}

	if _, ok := response["commit"]; !ok {
		t.Error("response should contain commit")
	}
}
