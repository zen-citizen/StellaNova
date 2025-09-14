package utils

import (
	"backend/internal/models"
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/paulmach/orb/geojson"
)

type mockGeoJsonManager struct {
	data map[string]interface{}
	err  error
}

func (m *mockGeoJsonManager) QueryPoint(ctx context.Context, lat, lng float64, city, layer string) (map[string]interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.data, nil
}

func (m *mockGeoJsonManager) GetAvailableLayers() map[string][]string {
	return map[string][]string{}
}

func TestExtractAttributesWithTransformer(t *testing.T) {
	tests := []struct {
		name         string
		transformer  PropertyTransformer
		mockData     map[string]interface{}
		mockError    error
		expectedData []models.Attribute
	}{
		{
			name:        "nil transformer uses default",
			transformer: nil,
			mockData: map[string]interface{}{
				"properties": geojson.Properties{
					"key1": "value1",
					"key2": "value2",
				},
			},
			expectedData: []models.Attribute{
				{Name: "key1", Value: "value1", IsFound: true},
				{Name: "key2", Value: "value2", IsFound: true},
			},
		},
		{
			name: "transformer success",
			transformer: func(props map[string]interface{}) ([]models.Attribute, error) {
				return []models.Attribute{
					{Name: "Custom", Value: "Value", IsFound: true},
				}, nil
			},
			mockData: map[string]interface{}{
				"properties": geojson.Properties{"any": "data"},
			},
			expectedData: []models.Attribute{
				{Name: "Custom", Value: "Value", IsFound: true},
			},
		},
		{
			name: "transformer error",
			transformer: func(props map[string]interface{}) ([]models.Attribute, error) {
				return nil, errors.New("transform failed")
			},
			mockData: map[string]interface{}{
				"properties": geojson.Properties{"any": "data"},
			},
			expectedData: []models.Attribute{},
		},
		{
			name:         "geomanager query error",
			transformer:  nil,
			mockError:    errors.New("geomanager error"),
			mockData:     nil,
			expectedData: []models.Attribute{},
		},
		{
			name:         "nil geomanager response",
			transformer:  nil,
			mockData:     nil,
			expectedData: []models.Attribute{},
		},
		{
			name:        "invalid properties type",
			transformer: nil,
			mockData: map[string]interface{}{
				"properties": "invalid",
			},
			expectedData: []models.Attribute{},
		},
		{
			name:        "nil values in properties filtered out",
			transformer: nil,
			mockData: map[string]interface{}{
				"properties": geojson.Properties{
					"key1": "value1",
					"key2": nil,
					"key3": "value3",
				},
			},
			expectedData: []models.Attribute{
				{Name: "key1", Value: "value1", IsFound: true},
				{Name: "key3", Value: "value3", IsFound: true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGeoManager := &mockGeoJsonManager{
				data: tt.mockData,
				err:  tt.mockError,
			}

			logger := slog.Default()
			ctx := context.Background()

			attrs := ExtractAttributes(ctx, mockGeoManager, 12.9, 77.6, "test-city", "test-layer", tt.transformer, logger)

			if len(attrs) != len(tt.expectedData) {
				t.Errorf("expected %d attributes, got %d", len(tt.expectedData), len(attrs))
			}

			for _, expected := range tt.expectedData {
				found := false
				for _, actual := range attrs {
					if actual.Name == expected.Name {
						if actual.Value != expected.Value {
							t.Errorf("attribute %q value: expected %q, got %q", expected.Name, expected.Value, actual.Value)
						}
						if actual.IsFound != expected.IsFound {
							t.Errorf("attribute %q IsFound: expected %v, got %v", expected.Name, expected.IsFound, actual.IsFound)
						}
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected attribute %q not found in results", expected.Name)
				}
			}

		})
	}
}

func TestBuildEntity(t *testing.T) {
	tests := []struct {
		name                    string
		attributes              []models.Attribute
		expectedAttributesCount int
		isAvailable             bool
	}{
		{
			name: "with found attributes",
			attributes: []models.Attribute{
				{Name: "Test", Value: "Value", IsFound: true},
			},
			expectedAttributesCount: 1,
			isAvailable:             true,
		},
		{
			name: "with not found attributes",
			attributes: []models.Attribute{
				{Name: "Test", Value: "Default", IsFound: false},
			},
			expectedAttributesCount: 0,
			isAvailable:             false,
		},
		{
			name:                    "with empty attributes",
			attributes:              []models.Attribute{},
			expectedAttributesCount: 0,
			isAvailable:             false,
		},
		{
			name: "mixed found and not found",
			attributes: []models.Attribute{
				{Name: "Found", Value: "Value", IsFound: true},
				{Name: "NotFound", Value: "Default", IsFound: false},
			},
			expectedAttributesCount: 2,
			isAvailable:             true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := BuildEntity(
				context.Background(),
				"Test Entity",
				"Not available message",
				nil,
				tt.attributes,
				slog.Default(),
			)

			if entity.IsAvailable != tt.isAvailable {
				t.Errorf("expected IsAvailable=%v, got %v", tt.isAvailable, entity.IsAvailable)
			}

			if len(entity.Attributes) != tt.expectedAttributesCount {
				t.Errorf("expected %d attributes, got %d", tt.expectedAttributesCount, len(entity.Attributes))
			}
		})
	}
}
