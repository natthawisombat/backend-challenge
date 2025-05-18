package handlers

import (
	"backend-challenge/entities"

	"github.com/gofiber/fiber/v2"
)

func Response(c *fiber.Ctx, response entities.Response, options ...map[string]interface{}) {
	// switch t := resp.(type) {
	// case render.ResponseBody:
	// 	utils.ToMap(resp, &response)
	// 	StatusCode = t.StatusCode
	// 	if t.Status != "OK" {
	// 		delete(response, "result")
	// 	}
	// }

	c.Status(response.StatusCode).JSON(response)
}
