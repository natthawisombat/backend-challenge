package routers

import (
	handlers "backend-challenge/adapters/http"
	mongo "backend-challenge/adapters/mongo"
	"backend-challenge/configs"
	"backend-challenge/entities"
	"backend-challenge/middlewares"
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
	auth.Post("/login", httpUser.Login)

	// //group protected with jwt
	users := prefix.Group("/users")
	users.Use(middlewares.JWTMiddleware())
	users.Get("/", httpUser.GetAll)
	users.Get("/:id", httpUser.Get)
	users.Patch("/:id", httpUser.Update)
	users.Delete("/:id", httpUser.Delete)

	prefix.Get("/healthcheck", func(c *fiber.Ctx) error {
		return handlers.Response(c, entities.Response{Status: "OK", Message: "Healthy"}, map[string]interface{}{"function": "Healthcheck"})
	})

	prefix.Use(func(c *fiber.Ctx) error {
		return handlers.Response(c, entities.Response{Status: "ER", ErrorCode: "ER404", ErrorMessage: "ไม่พบ Path", StatusCode: 404})
	})
}
