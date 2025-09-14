package cities

import (
	"backend/internal/utils"
	"context"
	"log/slog"
	"testing"
)

func TestRegistryInitialization(t *testing.T) {
	logger := slog.Default()
	geoManager := &utils.GeoJsonManager{}

	registry := NewRegistry(geoManager, logger)

	if registry == nil {
		t.Fatal("registry should not be nil")
	}

	cities := registry.SupportedCities()
	if len(cities) == 0 {
		t.Error("registry should have at least one city")
	}
}

func TestGetCityProvider(t *testing.T) {
	logger := slog.Default()
	geoManager := &utils.GeoJsonManager{}
	registry := NewRegistry(geoManager, logger)

	tests := []struct {
		name        string
		city        string
		shouldError bool
	}{
		{
			name:        "existing city",
			city:        "bengaluru",
			shouldError: false,
		},
		{
			name:        "non-existent city",
			city:        "unknown",
			shouldError: true,
		},
		{
			name:        "empty city",
			city:        "",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := registry.GetCityProvider(context.Background(), tt.city)

			if tt.shouldError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if provider == nil {
					t.Error("provider should not be nil")
				}
			}
		})
	}
}
