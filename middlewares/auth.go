package middlewares

import (
	"backend-challenge/entities"
	"backend-challenge/utils"
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or invalid token",
			})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: " + err.Error(),
			})
		}

		userID := claims["user_id"]
		ctx := context.WithValue(c.UserContext(), entities.UserIDKey, userID)
		c.SetUserContext(ctx)

		return c.Next()
	}
}
