package cities

import (
	"backend/internal/models"
	"backend/internal/utils"
	"context"
	"log/slog"
)

type ChennaiProvider struct {
	logger     *slog.Logger
	geoManager utils.GeoJsonManagerInterface
}

func NewChennaiProvider(geoManager utils.GeoJsonManagerInterface, logger *slog.Logger) *ChennaiProvider {
	return &ChennaiProvider{
		logger:     logger,
		geoManager: geoManager,
	}
}

func (p *ChennaiProvider) FormattedName() string {
	return "Chennai"
}

func (p *ChennaiProvider) Name() string {
	return "chennai"
}

func (p *ChennaiProvider) Bounds() *models.Bounds {
	return &models.Bounds{
		Northeast: models.Coordinate{Lat: 13.9, Lng: 80.19},
		Southwest: models.Coordinate{Lat: 12.9, Lng: 80.9},
	}
}

func (p *ChennaiProvider) GetEntities(ctx context.Context, lat, lng float64) ([]models.Entity, error) {
	p.logger.InfoContext(ctx, "getting entities", slog.String("city", p.Name()))

	var entities []models.Entity

	if entity := p.getGCCEntity(ctx, lat, lng); entity != nil {
		entities = append(entities, *entity)
	}

	return entities, nil
}

func (p *ChennaiProvider) getGCCEntity(ctx context.Context, lat, lng float64) *models.Entity {
	attributes := append(
		utils.ExtractAttributes(ctx, p.geoManager, lat, lng, p.Name(), "gcc", nil, p.logger),
	)

	entity := utils.BuildEntity(ctx, "Greater Chennai Corporation",
		"This information is unavailable for this address. This could be because the area is outside GBA Corporation limits.",
		nil,
		attributes, p.logger)

	return &entity
}
