package adapters

import (
	"backend-challenge/entities"
	"backend-challenge/pkg/logging"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Response(c *fiber.Ctx, response entities.Response, options ...map[string]interface{}) error {
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

	response.TransactionCode = fmt.Sprintf("%s", ctx.Value(entities.RequestId))
	statusCode := response.StatusCode
	response.StatusCode = 0
	return c.Status(statusCode).JSON(response)
}
