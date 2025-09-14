package services

import (
	"backend/internal/cities"
	"backend/internal/models"
	"context"
	"errors"
	"log/slog"
	"testing"
)

func TestGetEntitiesSuccess(t *testing.T) {
	mockRegistry := &mockCityRegistry{
		provider: &mockCityProvider{
			entities: []models.Entity{
				{Name: "Test Entity", IsAvailable: true},
			},
		},
	}

	service := NewEntitiesService(mockRegistry, slog.Default())

	req := &models.EntitiesRequest{
		Latitude:  12.9716,
		Longitude: 77.5946,
		City:      "bengaluru",
	}

	resp, err := service.GetEntities(context.Background(), req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Entities) != 1 {
		t.Errorf("expected 1 entity, got %d", len(resp.Entities))
	}
}

func TestGetEntitiesCityNotFound(t *testing.T) {
	mockRegistry := &mockCityRegistry{
		err: errors.New("unsupported city: unknown"),
	}

	service := NewEntitiesService(mockRegistry, slog.Default())

	req := &models.EntitiesRequest{
		Latitude:  12.9716,
		Longitude: 77.5946,
		City:      "unknown",
	}

	_, err := service.GetEntities(context.Background(), req)

	if err == nil {
		t.Error("expected error for unsupported city")
	}

	if err.Error() != "unsupported city: unknown" {
		t.Errorf("expected 'unsupported city: unknown', got '%s'", err.Error())
	}
}

func TestGetEntitiesProviderError(t *testing.T) {
	mockRegistry := &mockCityRegistry{
		provider: &mockCityProvider{
			err: errors.New("provider error"),
		},
	}

	service := NewEntitiesService(mockRegistry, slog.Default())

	req := &models.EntitiesRequest{
		Latitude:  12.9716,
		Longitude: 77.5946,
		City:      "bengaluru",
	}

	_, err := service.GetEntities(context.Background(), req)

	if err == nil {
		t.Error("expected provider error")
	}

	if err.Error() != "provider error" {
		t.Errorf("expected 'provider error', got '%s'", err.Error())
	}
}

type mockCityRegistry struct {
	provider cities.CityProvider
	err      error
}

func (m *mockCityRegistry) GetCityProvider(ctx context.Context, city string) (cities.CityProvider, error) {
	return m.provider, m.err
}

func (m *mockCityRegistry) SupportedCities() []string {
	return []string{"bengaluru"}
}

type mockCityProvider struct {
	entities []models.Entity
	err      error
}

func (m *mockCityProvider) Name() string           { return "mock" }
func (m *mockCityProvider) FormattedName() string  { return "Mock City" }
func (m *mockCityProvider) Bounds() *models.Bounds { return nil }
func (m *mockCityProvider) GetEntities(ctx context.Context, lat, lng float64) ([]models.Entity, error) {
	return m.entities, m.err
}
