package cities

import (
	"context"
	"fmt"
	"github.com/paulmach/orb/geojson"
	"log/slog"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAllCityProviders(t *testing.T) {
	logger := slog.Default()

	registry := NewRegistry(NewInstrumentedMockGeoJsonManager(), logger)
	registeredCities := registry.SupportedCities()

	if len(registeredCities) == 0 {
		t.Fatal("No city providers registered")
	}

	for _, city := range registeredCities {
		geoJsonPath := filepath.Join("..", "..", "assets", "geojson", city)
		if _, err := os.Stat(geoJsonPath); os.IsNotExist(err) {
			t.Errorf("Provider registered for '%s' but no GeoJSON directory at %s", city, geoJsonPath)
		}
	}

	for _, city := range registeredCities {
		t.Run(city, func(t *testing.T) {
			RunUniversalCityProviderTests(t, city)
		})
	}
}

func RunUniversalCityProviderTests(t *testing.T, cityName string) {
	ctx := context.Background()
	logger := slog.Default()

	expectedLayers := discoverCityLayers(t, cityName)
	if len(expectedLayers) == 0 {
		t.Fatalf("No GeoJSON layers found for city '%s'", cityName)
	}

	t.Logf("Discovered %d layers for %s: %v", len(expectedLayers), cityName, expectedLayers)

	mock := NewInstrumentedMockGeoJsonManager()
	mockRegistry := NewRegistry(mock, logger)
	provider, err := mockRegistry.GetCityProvider(ctx, cityName)
	if err != nil {
		t.Fatalf("Failed to get provider for '%s': %v", cityName, err)
	}

	t.Run("provider_properties", func(t *testing.T) {
		if provider.Name() != cityName {
			t.Errorf("Provider name mismatch: expected '%s', got '%s'", cityName, provider.Name())
		}

		if provider.FormattedName() == "" {
			t.Error("FormattedName() should not return empty string")
		}

		bounds := provider.Bounds()
		if bounds == nil {
			t.Error("Bounds() should not return nil")
		} else {
			if bounds.Northeast.Lat <= bounds.Southwest.Lat {
				t.Error("Northeast latitude should be greater than Southwest latitude")
			}
			if bounds.Northeast.Lng <= bounds.Southwest.Lng {
				t.Error("Northeast longitude should be greater than Southwest longitude")
			}
			if bounds.Northeast.Lat > 90 || bounds.Northeast.Lat < -90 {
				t.Error("Latitude must be between -90 and 90")
			}
			if bounds.Southwest.Lng > 180 || bounds.Southwest.Lng < -180 {
				t.Error("Longitude must be between -180 and 180")
			}
		}
	})

	t.Run("layer_discovery", func(t *testing.T) {
		mock.Reset()
		mock.SetMode(ModeValid)

		bounds := provider.Bounds()
		testLat := (bounds.Northeast.Lat + bounds.Southwest.Lat) / 2
		testLng := (bounds.Northeast.Lng + bounds.Southwest.Lng) / 2

		for _, layer := range expectedLayers {
			mock.AddResponse(cityName, layer, testLat, testLng, map[string]interface{}{
				"test_field": "test_value",
			})
		}

		entities, err := provider.GetEntities(ctx, testLat, testLng)
		if err != nil {
			t.Fatalf("GetEntities failed: %v", err)
		}

		if len(entities) == 0 {
			t.Error("No entities returned")
		}

		for _, layer := range expectedLayers {
			if !mock.WasLayerQueried(cityName, layer) {
				t.Errorf("Layer '%s' exists in GeoJSON but was not queried by provider", layer)
			}
		}

		t.Logf("Provider queried %d layers", len(mock.GetQueriedLayers(cityName)))
	})

	for _, layer := range expectedLayers {
		t.Run(fmt.Sprintf("layer_%s", layer), func(t *testing.T) {
			testLayerScenarios(t, provider, mock, cityName, layer)
		})
	}

	t.Run("combination_tests", func(t *testing.T) {
		bounds := provider.Bounds()
		testLat := (bounds.Northeast.Lat + bounds.Southwest.Lat) / 2
		testLng := (bounds.Northeast.Lng + bounds.Southwest.Lng) / 2

		t.Run("all_nil", func(t *testing.T) {
			mock.Reset()
			mock.SetMode(ModeNil)

			entities, err := provider.GetEntities(ctx, testLat, testLng)
			if err != nil {
				t.Fatalf("GetEntities failed: %v", err)
			}

			for _, entity := range entities {
				if entity.IsAvailable {
					t.Errorf("Entity '%s' should not be available when all layers return nil", entity.Name)
				}
				if entity.NotAvailableMessage == "" {
					t.Errorf("Entity '%s' has empty NotAvailableMessage when unavailable", entity.Name)
				}
			}
		})

		t.Run("mixed_valid_nil", func(t *testing.T) {
			mock.Reset()

			if len(expectedLayers) > 0 {
				mock.AddResponse(cityName, expectedLayers[0], testLat, testLng, map[string]interface{}{
					"test": "value",
				})
			}

			entities, err := provider.GetEntities(ctx, testLat, testLng)
			if err != nil {
				t.Fatalf("GetEntities failed: %v", err)
			}

			if len(entities) == 0 {
				t.Error("Should return entities even when some layers have no data")
			}
		})
	})
}

func testLayerScenarios(t *testing.T, provider CityProvider, mock *InstrumentedMockGeoJsonManager, cityName, layer string) {
	ctx := context.Background()
	bounds := provider.Bounds()
	testLat := (bounds.Northeast.Lat + bounds.Southwest.Lat) / 2
	testLng := (bounds.Northeast.Lng + bounds.Southwest.Lng) / 2

	t.Run("nil_response", func(t *testing.T) {
		mock.Reset()
		mock.SetMode(ModeNil)

		entities, err := provider.GetEntities(ctx, testLat, testLng)
		if err != nil {
			t.Fatalf("GetEntities failed with nil response: %v", err)
		}

		for _, entity := range entities {
			if entity.IsAvailable {
				t.Errorf("Entity '%s' should not be available with nil response", entity.Name)
			}
			if entity.NotAvailableMessage == "" {
				t.Errorf("Entity '%s' has empty NotAvailableMessage", entity.Name)
			}

			for _, attr := range entity.Attributes {
				if attr.IsFound {
					t.Errorf("Attribute '%s' should not be found with nil response", attr.Name)
				}
			}
		}
	})

	t.Run("valid_data", func(t *testing.T) {
		mock.Reset()
		mock.AddResponse(cityName, layer, testLat, testLng, map[string]interface{}{
			"name":   "Test Name",
			"value":  "Test Value",
			"number": 123,
		})

		entities, err := provider.GetEntities(ctx, testLat, testLng)
		if err != nil {
			t.Fatalf("GetEntities failed with valid data: %v", err)
		}

		foundAvailable := false
		for _, entity := range entities {
			if entity.IsAvailable {
				foundAvailable = true

				for _, attr := range entity.Attributes {
					if attr.IsFound && attr.Value == "" {
						t.Errorf("Attribute '%s' has IsFound=true but empty Value", attr.Name)
					}
				}
			}
		}

		if !foundAvailable && mock.WasLayerQueried(cityName, layer) {
			t.Errorf("Layer was queried but no entity became available - possibly not used by this provider")
		}
	})

	t.Run("wrong_types", func(t *testing.T) {
		mock.Reset()
		mock.AddResponse(cityName, layer, testLat, testLng, map[string]interface{}{
			"name":   123,
			"number": "not_number",
			"object": "not_object",
		})

		entities, err := provider.GetEntities(ctx, testLat, testLng)
		if err != nil {
			t.Fatalf("GetEntities failed with wrong types: %v", err)
		}

		if len(entities) == 0 {
			t.Error("Should return entities even with wrong types")
		}
	})

	t.Run("missing_fields", func(t *testing.T) {
		mock.Reset()
		mock.AddResponse(cityName, layer, testLat, testLng, map[string]interface{}{
			"partial_field": "exists",
		})

		entities, err := provider.GetEntities(ctx, testLat, testLng)
		if err != nil {
			t.Fatalf("GetEntities failed with missing fields: %v", err)
		}

		for _, entity := range entities {
			for _, attr := range entity.Attributes {
				if attr.IsFound && attr.Value == "" {
					t.Errorf("Attribute '%s' has IsFound=true but empty Value", attr.Name)
				}
				if !attr.IsFound && attr.Value == "" {
					t.Errorf("Attribute '%s' has IsFound=false and empty Value (should have default)", attr.Name)
				}
			}
		}
	})

	t.Run("null_values", func(t *testing.T) {
		mock.Reset()
		mock.AddResponse(cityName, layer, testLat, testLng, map[string]interface{}{
			"name":        nil,
			"valid_field": "valid",
		})

		entities, err := provider.GetEntities(ctx, testLat, testLng)
		if err != nil {
			t.Fatalf("GetEntities failed with null values: %v", err)
		}

		for _, entity := range entities {
			for _, attr := range entity.Attributes {
				if attr.Value == "" && attr.IsFound {
					t.Errorf("Attribute '%s' has empty Value with IsFound=true", attr.Name)
				}
			}
		}
	})

	t.Run("nested_issues", func(t *testing.T) {
		mock.Reset()
		mock.AddResponse(cityName, layer, testLat, testLng, map[string]interface{}{
			"nested": "not_an_object",
			"array":  "not_an_array",
			"valid": map[string]interface{}{
				"inner": "value",
			},
		})

		entities, err := provider.GetEntities(ctx, testLat, testLng)
		if err != nil {
			t.Fatalf("GetEntities failed with nested issues: %v", err)
		}

		if len(entities) == 0 {
			t.Error("Should return entities even with nested structure issues")
		}
	})
}

func discoverCityLayers(t *testing.T, cityName string) []string {
	geoJsonPath := filepath.Join("..", "..", "assets", "geojson", cityName)

	files, err := os.ReadDir(geoJsonPath)
	if err != nil {
		t.Logf("Warning: Could not read GeoJSON directory for %s: %v", cityName, err)
		return []string{}
	}

	var layers []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".geo.json") {
			layer := strings.TrimSuffix(file.Name(), ".geo.json")
			layers = append(layers, layer)
		}
	}

	return layers
}

type InstrumentedMockGeoJsonManager struct {
	queriedLayers map[string]map[string]int // city -> layer -> count
	responses     map[string]interface{}    // key -> response
	mode          ResponseMode
	logger        *slog.Logger
}

type ResponseMode int

const (
	ModeValid ResponseMode = iota
	ModeNil
	ModeMalformed
	ModeMixed
)

func NewInstrumentedMockGeoJsonManager() *InstrumentedMockGeoJsonManager {
	return &InstrumentedMockGeoJsonManager{
		queriedLayers: make(map[string]map[string]int),
		responses:     make(map[string]interface{}),
		mode:          ModeValid,
		logger:        slog.Default(),
	}
}

func (m *InstrumentedMockGeoJsonManager) Reset() {
	m.queriedLayers = make(map[string]map[string]int)
	m.responses = make(map[string]interface{})
	m.mode = ModeValid
}

func (m *InstrumentedMockGeoJsonManager) SetMode(mode ResponseMode) {
	m.mode = mode
}

func (m *InstrumentedMockGeoJsonManager) AddResponse(city, layer string, lat, lng float64, response map[string]interface{}) {
	key := fmt.Sprintf("%s/%s/%.4f,%.4f", city, layer, lat, lng)
	m.responses[key] = map[string]interface{}{
		"properties": geojson.Properties(response),
		"geometry":   "mock_geometry",
	}
}

func (m *InstrumentedMockGeoJsonManager) QueryPoint(ctx context.Context, lat, lng float64, city, layer string) (map[string]interface{}, error) {
	if m.queriedLayers[city] == nil {
		m.queriedLayers[city] = make(map[string]int)
	}
	m.queriedLayers[city][layer]++

	if m.mode == ModeNil {
		return nil, nil
	}

	key := fmt.Sprintf("%s/%s/%.4f,%.4f", city, layer, lat, lng)
	if response, exists := m.responses[key]; exists {
		if resp, ok := response.(map[string]interface{}); ok {
			return resp, nil
		}
	}

	switch m.mode {
	case ModeMalformed:
		return m.generateMalformedResponse(), nil
	case ModeMixed:
		if rand.Float32() > 0.5 {
			return nil, nil
		}
		return m.generateValidResponse(), nil
	case ModeValid:
		return m.generateValidResponse(), nil
	default:
		return nil, nil
	}
}

func (m *InstrumentedMockGeoJsonManager) GetAvailableLayers() map[string][]string {
	layers := make(map[string][]string)
	for city, cityLayers := range m.queriedLayers {
		for layer := range cityLayers {
			layers[city] = append(layers[city], layer)
		}
	}
	return layers
}

func (m *InstrumentedMockGeoJsonManager) WasLayerQueried(city, layer string) bool {
	if cityLayers, exists := m.queriedLayers[city]; exists {
		return cityLayers[layer] > 0
	}
	return false
}

func (m *InstrumentedMockGeoJsonManager) GetQueriedLayers(city string) []string {
	var layers []string
	if cityLayers, exists := m.queriedLayers[city]; exists {
		for layer := range cityLayers {
			layers = append(layers, layer)
		}
	}
	return layers
}

func (m *InstrumentedMockGeoJsonManager) generateValidResponse() map[string]interface{} {
	return map[string]interface{}{
		"properties": geojson.Properties{
			"name":   "Test Name",
			"value":  "Test Value",
			"number": 123,
		},
		"geometry": "mock_geometry",
	}
}

func (m *InstrumentedMockGeoJsonManager) generateMalformedResponse() map[string]interface{} {
	return map[string]interface{}{
		"properties": geojson.Properties{
			"name":   123,
			"nested": "not_object",
			"nil":    nil,
		},
		"geometry": "mock_geometry",
	}
}
