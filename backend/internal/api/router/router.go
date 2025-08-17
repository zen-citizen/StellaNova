package router

import (
	"backend/internal/api/handlers"
	"backend/internal/api/middleware"
	"backend/internal/cities"
	"backend/internal/config"
	"backend/internal/services"
	"backend/internal/utils"
	"backend/internal/version"
	"github.com/gin-gonic/gin"
	"log/slog"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func Setup(cfg *config.Config, geoManager *utils.GeoJsonManager, logger *slog.Logger) *gin.Engine {
	gin.SetMode(cfg.GinMode)

	r := gin.New()
	r.Use(otelgin.Middleware("StellaNova-backend"))
	r.Use(middleware.Logger(logger))

	logger.Info("setting up http router")

	cityRegistry := cities.NewRegistry(geoManager, logger)

	entitiesService := services.NewEntitiesService(cityRegistry, logger)

	entitiesHandler := handlers.NewEntitiesHandler(entitiesService, logger)

	r.GET("/health", func(c *gin.Context) {
		logger.DebugContext(c.Request.Context(), "health check accessed")
		c.JSON(200, gin.H{
			"status":     "up",
			"build_time": version.GetBuildTime(),
			"commit":     version.GetCommitSHA(),
			"dependencies": gin.H{
				"geojson": geoManager.GetAvailableLayers(),
			},
		})
	})

	api := r.Group("/api/v1")
	{
		api.GET("/entities", entitiesHandler.GetEntities)
	}

	logger.Info("router setup")
	return r
}
