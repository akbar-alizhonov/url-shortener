package middlewares

import (
	"awesomeProject/pkg/logger"
	"log/slog"
	"time"

	"github.com/labstack/echo/v5"
)

func RequestLogger(log *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			start := time.Now()
			err := next(c)
			req := c.Request()
			reqID := logger.RequestIDFromContext(req.Context())
			log.Info("http-server",
				"time", start,
				"method", req.Method,
				"url", req.URL.String(),
				"request_id", reqID,
			)
			return err
		}
	}
}
