package middlewares

import (
	"awesomeProject/pkg/logger"

	"github.com/labstack/echo/v5"
)

func RequestContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		reqID := c.Response().Header().Get(echo.HeaderXRequestID)
		ctx := logger.WithRequestID(c.Request().Context(), reqID)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
