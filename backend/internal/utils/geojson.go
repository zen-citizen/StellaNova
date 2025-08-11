package utils

import (
	"context"
	"fmt"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/planar"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type GeoJsonManager struct {
	data   map[string]*geojson.FeatureCollection
	logger *slog.Logger
}

func NewGeoJsonManager(logger *slog.Logger) (*GeoJsonManager, error) {
	manager := &GeoJsonManager{
		data:   make(map[string]*geojson.FeatureCollection),
		logger: logger,
	}

	err := manager.loadAllGeoJson("assets")
	if err != nil {
		return nil, fmt.Errorf("failed to load geojson data: %w", err)
	}

	return manager, nil
}

func (g *GeoJsonManager) loadAllGeoJson(assetsPath string) error {
	geoJsonPath := filepath.Join(assetsPath, "geojson")

	g.logger.Info("loading geojson data into memory",
		slog.String("path", geoJsonPath),
	)

	totalFiles := 0
	totalSize := int64(0)

	cities, err := os.ReadDir(geoJsonPath)
	if err != nil {
		return fmt.Errorf("failed to read geojson directoy: %w", err)
	} else {
		g.logger.Info("found cities",
			slog.Any("cities", cities),
		)
	}

	for _, cityDir := range cities {
		if !cityDir.IsDir() {
			continue
		}

		cityName := cityDir.Name()
		cityPath := filepath.Join(geoJsonPath, cityName)

		files, err := os.ReadDir(cityPath)
		if err != nil {
			g.logger.Warn("failed to read city directory",
				slog.String("city", cityName),
				slog.Any("error", err),
			)
			continue
		}

		cityFiles := 0
		for _, file := range files {
			if !strings.HasSuffix(file.Name(), ".geo.json") {
				continue
			}

			layerName := strings.TrimSuffix(file.Name(), ".geo.json")
			filePath := filepath.Join(cityPath, file.Name())

			fileInfo, err := file.Info()
			if err != nil {
				g.logger.Warn("failed to get file info",
					slog.String("file", file.Name()),
					slog.String("city", cityName),
					slog.Any("error", err),
				)
				continue
			}

			err = g.loadGeoJsonFile(cityName, layerName, filePath)
			if err != nil {
				g.logger.Warn("Failed to load GeoJSON file",
					slog.String("file", filePath),
					slog.Any("error", err),
				)
				continue
			}

			totalFiles++
			cityFiles++
			totalSize += fileInfo.Size()
		}

		g.logger.Debug("Loaded city GeoJSON data",
			slog.String("city", cityName),
			slog.Int("files", cityFiles),
		)
	}

	g.logger.Info("GeoJSON data loaded successfully",
		slog.Int("total_files", totalFiles),
		slog.Int("cities", len(cities)),
		slog.String("total_size", formatBytes(totalSize)),
		slog.Int("layers_in_memory", len(g.data)),
	)

	return nil
}

func (g *GeoJsonManager) loadGeoJsonFile(city, layer, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	fc, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		return fmt.Errorf("failed to parse geojson: %w", err)
	}

	key := fmt.Sprintf("%s/%s", city, layer)
	g.data[key] = fc

	g.logger.Debug("Loaded GeoJSON layer",
		slog.String("city", city),
		slog.String("layer", layer),
		slog.Int("features", len(fc.Features)),
		slog.String("size", formatBytes(int64(len(data)))),
	)

	return nil
}

func (g *GeoJsonManager) QueryPoint(ctx context.Context, lat, lng float64, city, layer string) (map[string]interface{}, error) {
	userPoint := orb.Point{lng, lat}

	g.logger.DebugContext(ctx, "Querying GeoJSON for point",
		slog.Float64("latitude", lat),
		slog.Float64("longitude", lng),
		slog.String("city", city),
		slog.String("layer", layer),
	)

	key := fmt.Sprintf("%s/%s", city, layer)
	fc, exists := g.data[key]
	if !exists {
		g.logger.DebugContext(ctx, "GeoJSON layer not found",
			slog.String("key", key),
		)
		return nil, fmt.Errorf("geojson layer not found: %s", key)
	}

	for _, feature := range fc.Features {
		if g.pointInFeature(userPoint, feature) {
			g.logger.DebugContext(ctx, "Found matching feature",
				slog.String("layer", layer),
				slog.String("city", city),
			)
			return map[string]interface{}{
				"properties": feature.Properties,
				"geometry":   feature.Geometry,
			}, nil
		}
	}

	g.logger.DebugContext(ctx, "No matching feature found",
		slog.String("layer", layer),
		slog.String("city", city),
	)
	return nil, nil
}

func (g *GeoJsonManager) pointInFeature(point orb.Point, feature *geojson.Feature) bool {
	switch geom := feature.Geometry.(type) {
	case orb.Polygon:
		return planar.PolygonContains(geom, point)
	case orb.MultiPolygon:
		return planar.MultiPolygonContains(geom, point)
	default:
		return false
	}
}

func (g *GeoJsonManager) GetAvailableLayers() map[string][]string {
	layers := make(map[string][]string)

	for key := range g.data {
		parts := strings.Split(key, "/")
		if len(parts) == 2 {
			city, layer := parts[0], parts[1]
			layers[city] = append(layers[city], layer)
		}
	}

	return layers
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
