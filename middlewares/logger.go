package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func ZapLoggerMiddleware(logger *zap.SugaredLogger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		stop := time.Now()

		latency := stop.Sub(start)

		logger.Infow("HTTP Request",
			"method", c.Method(),
			"path", c.OriginalURL(),
			"status", c.Response().StatusCode(),
			"latency", latency.String(),
			"ip", c.IP(),
			"user_agent", c.Get("User-Agent"),
		)

		return err
	}
}
