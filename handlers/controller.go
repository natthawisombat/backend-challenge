package handlers

import "github.com/gofiber/fiber/v2"

func WrapHandler(f func(*fiber.Ctx)) fiber.Handler {
	return func(c *fiber.Ctx) error {
		f(c)
		return nil
	}
}
