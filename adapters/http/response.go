package adapters

import (
	"backend-challenge/entities"
	"backend-challenge/pkg/logging"

	"github.com/gofiber/fiber/v2"
)

func Response(c *fiber.Ctx, response entities.Response, options ...map[string]interface{}) {
	ctx := c.UserContext()
	logger := logging.FromContext(ctx)

	fields := []interface{}{
		"status_code", response.StatusCode,
		"status", response.Status,
	}

	if len(options) > 0 {
		for k, v := range options[0] {
			fields = append(fields, k, v)
		}
	}

	if response.Status == "OK" {
		logger.Infow("response success", fields...)
	} else {
		fields = append(fields, "error_message", response.ErrorMessage)
		logger.Errorw("response error", fields...)
	}

	c.Status(response.StatusCode).JSON(response)
}
