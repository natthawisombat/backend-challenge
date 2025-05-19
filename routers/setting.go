package routers

import (
	handlers "backend-challenge/adapters/http"
	mongo "backend-challenge/adapters/mongo"
	"backend-challenge/configs"
	"backend-challenge/entities"
	"backend-challenge/pkg/logging"
	"backend-challenge/usecases"
	"fmt"
	"time"

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
	auth.Post("/register", handlers.WrapHandler(httpUser.Create))
	// auth.Post("/login")

	// //group protected with jwt
	// users := prefix.Group("/users")
	// users.Get("/")
	// users.Get("/:id")
	// users.Patch("/:id")
	// users.Delete("/:id")

	prefix.Use(handlers.WrapHandler(func(c *fiber.Ctx) {
		ctx := c.UserContext()
		logger := logging.FromContext(ctx)
		fmt.Println(ctx.Value(entities.RequestId))
		coll := cfg.DBMongo.DB.Collection("user")
		newUser := entities.User{Name: "Tee", Email: "tee@hotmail.com", Password: "Hash", CreatedAt: time.Now()}
		result, err := coll.InsertOne(ctx, newUser)
		if err != nil {
			handlers.Response(c, entities.Response{Status: "ER", ErrorCode: "ER999", ErrorMessage: err.Error(), StatusCode: 500})
		}
		fmt.Println(result)
		time.Sleep(10 * time.Second)
		logger.Infow("insert user success", "inserted_id", result.InsertedID)
		// handlers.Response(c, entities.Response{Status: "ER", ErrorCode: "ER404", ErrorMessage: "ไม่พบ Path", StatusCode: 404})
		handlers.Response(c, entities.Response{Status: "OK", StatusCode: 200})
	}))
}
