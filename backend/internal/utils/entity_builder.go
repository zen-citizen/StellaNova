package utils

import (
	"backend/internal/models"
	"context"
	"fmt"
	"github.com/paulmach/orb/geojson"
	"log/slog"
	"reflect"
)

type PropertyTransformer func(map[string]interface{}) ([]models.Attribute, error)

func ExtractAttributes(ctx context.Context, geoManager *GeoJsonManager, lat, lng float64, city, layer string, transformer PropertyTransformer, logger *slog.Logger) []models.Attribute {
	geoData, err := geoManager.QueryPoint(ctx, lat, lng, city, layer)
	if err != nil {
		logger.ErrorContext(ctx, "failed to query geojson data",
			slog.String("city", city),
			slog.String("layer", layer),
			slog.Any("error", err),
		)
		return []models.Attribute{}
	}

	if geoData == nil {
		return []models.Attribute{}
	}

	properties, ok := geoData["properties"].(geojson.Properties)
	if !ok {
		logger.ErrorContext(ctx, "could not read properties",
			slog.Any("properties", geoData["properties"]),
			slog.Any("type", reflect.TypeOf(geoData["properties"])),
			slog.String("city", city),
			slog.String("layer", layer),
			slog.Float64("lat", lat),
			slog.Float64("lng", lng),
		)
		return []models.Attribute{}
	}

	var attributes []models.Attribute
	if transformer != nil {
		attributes, err = transformer(properties)
		if err != nil {
			logger.WarnContext(ctx, "Failed to transform properties",
				slog.String("city", city),
				slog.String("layer", layer),
				slog.Any("error", err),
			)
		}
	} else {
		for key, value := range properties {
			if value != nil {
				attributes = append(attributes, models.Attribute{
					Name:    key,
					Value:   fmt.Sprintf("%v", value),
					IsFound: true,
				})
			}
		}
	}

	return attributes
}

func BuildEntity(ctx context.Context, name, notAvailableMsg string, disclaimer *string, attributes []models.Attribute, logger *slog.Logger) models.Entity {
	logger.DebugContext(ctx, "entity built successfully",
		slog.String("entity", name),
		slog.Int("attribute_count", len(attributes)),
	)

	allAttributesNotFound := true
	for _, attribute := range attributes {
		if attribute.IsFound {
			allAttributesNotFound = false
			break
		}
	}
	if len(attributes) == 0 || allAttributesNotFound {
		return models.NewUnavailableEntity(name, notAvailableMsg, disclaimer)
	}

	return models.NewAvailableEntity(name, notAvailableMsg, disclaimer, attributes)
}
