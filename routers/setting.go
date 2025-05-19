package routers

import (
	handlers "backend-challenge/adapters/http"
	mongo "backend-challenge/adapters/mongo"
	"backend-challenge/configs"
	"backend-challenge/entities"
	"backend-challenge/usecases"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(cfg *configs.Setting) {
	validate := validator.New()
	prefix := cfg.App.Group(configs.App.Prefix)

	repository := mongo.NewMongoRepository(cfg.DBMongo.DB)
	httpUser := usecases.NewHttpUser(validate, repository)
	//group auth
	auth := prefix.Group("/auth")
	auth.Post("/register", httpUser.Create)
	// auth.Post("/login")

	// //group protected with jwt
	// users := prefix.Group("/users")
	// users.Get("/")
	// users.Get("/:id")
	// users.Patch("/:id")
	// users.Delete("/:id")

	prefix.Use(func(c *fiber.Ctx) error {
		return handlers.Response(c, entities.Response{Status: "ER", ErrorCode: "ER404", ErrorMessage: "ไม่พบ Path", StatusCode: 404})
	})
}
