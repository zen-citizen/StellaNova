package cities

import (
	"backend/internal/utils"
	"fmt"
	"log/slog"
)

type Registry struct {
	cities     map[string]CityProvider
	geoManager *utils.GeoJsonManager
	logger     *slog.Logger
}

func NewRegistry(geoManager *utils.GeoJsonManager, logger *slog.Logger) *Registry {
	r := &Registry{
		cities:     make(map[string]CityProvider),
		logger:     logger,
		geoManager: geoManager,
	}

	r.cities["bengaluru"] = NewBangaloreProvider(geoManager, logger)

	r.logger.Info("City registry initialized",
		slog.Int("registered_cities", len(r.cities)),
		slog.Any("cities", r.SupportedCities()),
		slog.Any("available_layers", geoManager.GetAvailableLayers()),
	)

	return r
}

func (r *Registry) GetCityProvider(city string) (CityProvider, error) {
	provider, exists := r.cities[city]
	if !exists {
		r.logger.Warn("unsupported city requested",
			slog.String("requested_city", city),
			slog.Any("supported_cities", r.SupportedCities()),
		)
		return nil, fmt.Errorf("unsupported city: %s", city)
	}

	r.logger.Debug("City provider retrieved", slog.String("city", city))
	return provider, nil
}

func (r *Registry) SupportedCities() []string {
	cities := make([]string, 0, len(r.cities))
	for city := range r.cities {
		cities = append(cities, city)
	}
	return cities
}

func (r *Registry) GetGeoJSONManager() *utils.GeoJsonManager {
	return r.geoManager
}
