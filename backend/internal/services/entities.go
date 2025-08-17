package services

import (
	"backend/internal/cities"
	"backend/internal/models"
	"context"
	"log/slog"
)

type EntitiesService struct {
	cityRegistry *cities.Registry
	logger       *slog.Logger
}

func NewEntitiesService(cityRegistry *cities.Registry, logger *slog.Logger) *EntitiesService {
	logger.Info("entities service initialized")
	return &EntitiesService{
		cityRegistry: cityRegistry,
		logger:       logger,
	}
}

func (s *EntitiesService) GetEntities(ctx context.Context, req *models.EntitiesRequest) (*models.EntitiesResponse, error) {
	s.logger.InfoContext(ctx, "processing entities request",
		slog.Float64("latitude", req.Latitude),
		slog.Float64("longitude", req.Longitude),
		slog.String("city", req.City),
	)

	cityProvider, err := s.cityRegistry.GetCityProvider(ctx, req.City)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to get city provider",
			slog.String("city", req.City),
			slog.Any("error", err),
		)
		return nil, err
	}

	entities, err := cityProvider.GetEntities(ctx, req.Latitude, req.Longitude)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to get entities from city provider",
			slog.String("city", req.City),
			slog.Any("error", err),
		)
		return nil, err
	}

	s.logger.InfoContext(ctx, "location request completed",
		slog.String("city", req.City),
		slog.Int("entity_count", len(entities)),
	)

	return &models.EntitiesResponse{
		Entities: entities,
	}, nil
}
