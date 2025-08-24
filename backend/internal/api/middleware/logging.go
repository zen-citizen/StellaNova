package middleware

import (
	"github.com/gin-gonic/gin"
	"log/slog"
)

func Logger(logger *slog.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.InfoContext(param.Request.Context(),
			"http request",
			slog.String("method", param.Method),
			slog.String("path", param.Path),
			slog.Int("status", param.StatusCode),
			slog.Int("status", param.StatusCode),
			slog.Time("timestamp", param.TimeStamp),
			slog.Any("query_params", param.Request.URL.Query()),
		)
		return ""
	})
}
