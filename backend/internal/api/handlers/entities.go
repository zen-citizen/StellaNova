package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EntitiesHandler struct {
	service *services.EntitiesService
	logger  *slog.Logger
}

func NewEntitiesHandler(service *services.EntitiesService, logger *slog.Logger) *EntitiesHandler {
	return &EntitiesHandler{
		service: service,
		logger:  logger,
	}
}

func (h *EntitiesHandler) GetEntities(c *gin.Context) {
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	city := c.Query("city")

	if latStr == "" || lngStr == "" {
		h.logger.Warn("invalid location request - missing parameters",
			slog.String("lat", latStr),
			slog.String("lng", lngStr),
			slog.String("city", city),
			slog.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "lat, lng, and city parameters are required",
		})
		return
	}

	if city == "" {
		h.logger.Info("city query param empty, using bengaluru")
		city = "bengaluru"
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		h.logger.Warn("invalid latitude format",
			slog.String("lat", latStr),
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid latitude",
		})
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		h.logger.Warn("invalid longitude format",
			slog.String("lng", lngStr),
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid longitude",
		})
		return
	}

	req := &models.EntitiesRequest{
		Latitude:  lat,
		Longitude: lng,
		City:      city,
	}

	h.logger.Info("processing location request",
		slog.Float64("latitude", lat),
		slog.Float64("longitude", lng),
		slog.String("city", city),
		slog.String("client_ip", c.ClientIP()),
	)

	resp, err := h.service.GetEntities(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("location request failed",
			slog.Float64("latitude", lat),
			slog.Float64("longitude", lng),
			slog.String("city", city),
			slog.Any("error", err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("Location request completed successfully",
		slog.Float64("latitude", lat),
		slog.Float64("longitude", lng),
		slog.String("city", city),
		slog.Int("entity_count", len(resp.Entities)),
	)

	c.JSON(http.StatusOK, resp)
}
