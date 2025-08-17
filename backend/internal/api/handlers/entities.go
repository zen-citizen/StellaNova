package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const MaxLatitude = 90.0
const MinLatitude = -90.0
const MaxLongitude = 180.0
const MinLongitude = -180.0

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
		h.logger.WarnContext(c.Request.Context(),
			"invalid entities request - missing parameters",
			slog.String("lat", latStr),
			slog.String("lng", lngStr),
			slog.String("city", city),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "lat, lng, and city parameters are required",
		})
		return
	}

	if city == "" {
		h.logger.InfoContext(c.Request.Context(), "city query param empty, using bengaluru")
		city = "bengaluru"
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil || lat < MinLatitude || lat > MaxLatitude {
		h.logger.WarnContext(c.Request.Context(),
			"invalid latitude format",
			slog.String("lat", latStr),
			slog.Float64("min_lat", MinLatitude),
			slog.Float64("max_lat", MaxLatitude),
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid latitude",
		})
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil || lng < MinLongitude || lng > MaxLongitude {
		h.logger.WarnContext(c.Request.Context(),
			"invalid longitude format",
			slog.String("lng", lngStr),
			slog.Float64("min_lng", MinLongitude),
			slog.Float64("max_lng", MaxLongitude),
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

	h.logger.InfoContext(c.Request.Context(),
		"processing entities request",
		slog.Float64("latitude", lat),
		slog.Float64("longitude", lng),
		slog.String("city", city),
	)

	resp, err := h.service.GetEntities(c.Request.Context(), req)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(),
			"entities request failed",
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

	h.logger.InfoContext(c.Request.Context(),
		"entities request completed successfully",
		slog.Float64("latitude", lat),
		slog.Float64("longitude", lng),
		slog.String("city", city),
		slog.Int("entity_count", len(resp.Entities)),
	)

	c.JSON(http.StatusOK, resp)
}
