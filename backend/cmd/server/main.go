package main

import (
	"backend/internal/api/router"
	"backend/internal/config"
	"backend/internal/utils"
	"context"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("error loading .env file")
	}

	cfg := config.New()

	logger := utils.SetupLogger(cfg.LogLevel, cfg.LogFormat)
	slog.SetDefault(logger)

	ctx := context.Background()
	logger.InfoContext(ctx, "starting application",
		slog.String("port", cfg.Port),
		slog.String("gin_mode", cfg.GinMode),
		slog.String("log_level", cfg.LogLevel),
	)

	geoManager, err := utils.NewGeoJsonManager(logger)
	if err != nil {
		logger.ErrorContext(ctx, "failed to initialize geojson manager", slog.Any("error", err))
		os.Exit(1)
	}

	r := router.Setup(cfg, geoManager, logger)

	logger.InfoContext(ctx, "starting server", slog.String("port", cfg.Port))
	if err := r.Run(":" + cfg.Port); err != nil {
		logger.ErrorContext(ctx, "failed to start server", slog.Any("error", err))
		os.Exit(1)
	}
}
