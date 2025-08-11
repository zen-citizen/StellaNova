package router

import (
	"backend/internal/api/handlers"
	"backend/internal/api/middleware"
	"backend/internal/cities"
	"backend/internal/config"
	"backend/internal/services"
	"backend/internal/utils"
	"github.com/gin-gonic/gin"
	"log/slog"
)

func Setup(cfg *config.Config, geoManager *utils.GeoJsonManager, logger *slog.Logger) *gin.Engine {
	gin.SetMode(cfg.GinMode)

	r := gin.New()
	r.Use(middleware.Logger(logger))

	logger.Info("setting up http router")

	cityRegistry := cities.NewRegistry(geoManager, logger)

	entitiesService := services.NewEntitiesService(cityRegistry, logger)

	entitiesHandler := handlers.NewEntitiesHandler(entitiesService, logger)

	api := r.Group("/api/v1")
	{
		api.GET("/entities", entitiesHandler.GetEntities)

		api.GET("/health", func(c *gin.Context) {
			logger.Debug("health check accessed")
			c.JSON(200, gin.H{
				"status": "up",
			})
		})
	}

	logger.Info("router setup")
	return r
}
