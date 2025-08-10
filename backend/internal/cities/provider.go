package cities

import (
	"backend/internal/models"
	"context"
)

type CityProvider interface {
	Name() string
	FormattedName() string
	Bounds() *models.Bounds

	GetEntities(ctx context.Context, lat, lng float64) ([]models.Entity, error)
}
