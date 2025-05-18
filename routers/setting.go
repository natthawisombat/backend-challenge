package routers

import (
	"backend-challenge/configs"
	"backend-challenge/entities"
	"backend-challenge/handlers"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(cfg *configs.Setting) {
	// setup routes prefix
	prefix := cfg.App.Group(configs.App.Prefix)

	prefix.Use(handlers.WrapHandler(func(c *fiber.Ctx) {

		coll := cfg.DBMongo.DB.Collection("user")
		newUser := entities.User{Name: "Tee", Email: "tee@hotmail.com", Password: "Hash", CreatedAt: time.Now()}
		result, err := coll.InsertOne(c.Context(), newUser)
		if err != nil {
			handlers.Response(c, entities.Response{Status: "ER", ErrorCode: "ER999", ErrorMessage: err.Error(), StatusCode: 500})
		}
		fmt.Println(result)
		time.Sleep(10 * time.Second)
		// handlers.Response(c, entities.Response{Status: "ER", ErrorCode: "ER404", ErrorMessage: "ไม่พบ Path", StatusCode: 404})
		handlers.Response(c, entities.Response{Status: "OK", StatusCode: 200})
	}))
}
