package middlewares

import (
	handlers "backend-challenge/adapters/http"
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
			return handlers.Response(c, entities.Response{Status: "ER", ErrorMessage: "Missing or invalid token", ErrorCode: "ER401", StatusCode: 401})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			return handlers.Response(c, entities.Response{Status: "ER", ErrorMessage: "Unauthorized: " + err.Error(), ErrorCode: "ER401", StatusCode: 401})
		}

		userID := claims["user_id"]
		ctx := context.WithValue(c.UserContext(), entities.UserIDKey, userID)
		c.SetUserContext(ctx)

		return c.Next()
	}
}
