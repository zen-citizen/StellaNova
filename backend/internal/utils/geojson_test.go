package utils

import (
	"context"
	"github.com/paulmach/orb/geojson"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
)

func TestGeoJsonFileLoading(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		shouldError bool
	}{
		{
			name: "valid geojson",
			content: `{
				"type": "FeatureCollection",
				"features": []
			}`,
			shouldError: false,
		},
		{
			name:        "invalid json",
			content:     `{invalid json}`,
			shouldError: true,
		},
		{
			name:        "not a feature collection",
			content:     `{"type": "Point", "coordinates": [0, 0]}`,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "test.geo.json")
			writeFileErr := os.WriteFile(testFile, []byte(tt.content), 0644)
			if writeFileErr != nil {
				t.Fatalf("failed to create temp file: %v", writeFileErr)
			}

			manager := &GeoJsonManager{
				data:   make(map[string]*geojson.FeatureCollection),
				logger: slog.Default(),
			}

			err := manager.loadGeoJsonFile("test", "layer", testFile)

			if tt.shouldError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestPointInPolygonQueries(t *testing.T) {
	tests := []struct {
		name     string
		lat      float64
		lng      float64
		expected bool
	}{
		{
			name:     "point inside polygon",
			lat:      0.5,
			lng:      0.5,
			expected: true,
		},
		{
			name:     "point outside polygon",
			lat:      2.0,
			lng:      2.0,
			expected: false,
		},
		{
			name:     "point on boundary",
			lat:      0.0,
			lng:      0.0,
			expected: true,
		},
	}

	geoJSON := `{
		"type": "FeatureCollection",
		"features": [{
			"type": "Feature",
			"properties": {"name": "test"},
			"geometry": {
				"type": "Polygon",
				"coordinates": [[[0,0], [1,0], [1,1], [0,1], [0,0]]]
			}
		}]
	}`

	tmpDir := t.TempDir()
	err := os.MkdirAll(filepath.Join(tmpDir, "test_city"), 0755)
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	testFile := filepath.Join(tmpDir, "test_city", "test.geo.json")
	err = os.WriteFile(testFile, []byte(geoJSON), 0644)
	if err != nil {
		t.Fatalf("failed to write geojson file: %v", err)
	}

	manager := &GeoJsonManager{
		data:   make(map[string]*geojson.FeatureCollection),
		logger: slog.Default(),
	}
	err = manager.loadGeoJsonFile("test_city", "test", testFile)
	if err != nil {
		t.Fatalf("failed to load geojson file: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := manager.QueryPoint(context.Background(), tt.lat, tt.lng, "test_city", "test")

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.expected && result == nil {
				t.Error("expected to find point in polygon")
			}
			if !tt.expected && result != nil {
				t.Error("expected point to be outside polygon")
			}
		})
	}
}

func TestMultiPolygonSupport(t *testing.T) {
	geoJSON := `{
		"type": "FeatureCollection",
		"features": [{
			"type": "Feature",
			"properties": {"name": "multi"},
			"geometry": {
				"type": "MultiPolygon",
				"coordinates": [
					[[[0,0], [1,0], [1,1], [0,1], [0,0]]],
					[[[2,2], [3,2], [3,3], [2,3], [2,2]]]
				]
			}
		}]
	}`

	tmpDir := t.TempDir()
	err := os.MkdirAll(filepath.Join(tmpDir, "test_city"), 0755)
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	testFile := filepath.Join(tmpDir, "test_city", "multi.geo.json")
	err = os.WriteFile(testFile, []byte(geoJSON), 0644)
	if err != nil {
		t.Fatalf("failed to write geojson file: %v", err)
	}

	manager := &GeoJsonManager{
		data:   make(map[string]*geojson.FeatureCollection),
		logger: slog.Default(),
	}
	err = manager.loadGeoJsonFile("test_city", "multi", testFile)
	if err != nil {
		t.Fatalf("failed to load geojson file: %v", err)
	}

	tests := []struct {
		name     string
		lat      float64
		lng      float64
		expected bool
	}{
		{"in first polygon", 0.5, 0.5, true},
		{"in second polygon", 2.5, 2.5, true},
		{"outside both", 5.0, 5.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := manager.QueryPoint(context.Background(), tt.lat, tt.lng, "test_city", "multi")

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.expected && result == nil {
				t.Error("expected to find point in multipolygon")
			}
			if !tt.expected && result != nil {
				t.Error("expected point to be outside multipolygon")
			}
		})
	}
}
