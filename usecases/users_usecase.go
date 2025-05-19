package usecases

import (
	handlers "backend-challenge/adapters/http"
	"backend-challenge/entities"
	"backend-challenge/utils"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type HttpUser struct {
	repo     userRepository
	validate *validator.Validate
}

func NewHttpUser(validate *validator.Validate, repo userRepository) HttpUser {
	return HttpUser{validate: validate, repo: repo}
}

func (uc *HttpUser) Create(c *fiber.Ctx) {
	var bodyRequest entities.User
	if err := c.BodyParser(&bodyRequest); err != nil {
		handlers.Response(c, entities.Response{Status: "ER", ErrorMessage: err.Error(), ErrorCode: "ER400", StatusCode: 400}, map[string]interface{}{"function": "Create"})
		return
	}

	err := uc.validate.Struct(bodyRequest)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			handlers.Response(c, entities.Response{Status: "ER", ErrorMessage: err.Error(), ErrorCode: "ER400", StatusCode: 400}, map[string]interface{}{"function": "Create"})
			return
		}
	}

	// Todo
	// check duplicate
	if err := uc.repo.CheckDuplicateUser(bodyRequest.Email, c.UserContext()); err != nil {
		handlers.Response(c, entities.Response{Status: "ER", ErrorMessage: err.Error(), ErrorCode: "ER400", StatusCode: 400}, map[string]interface{}{"function": "Create"})
		return
	}
	// hash password
	bodyRequest.Password = utils.Hash(bodyRequest.Password)
	bodyRequest.CreatedAt = time.Now()
	if err := uc.repo.Register(bodyRequest, c.UserContext()); err != nil {
		handlers.Response(c, entities.Response{Status: "ER", ErrorMessage: err.Error(), ErrorCode: "ER400", StatusCode: 400}, map[string]interface{}{"function": "Create"})
		return
	}

	handlers.Response(c, entities.Response{Status: "OK", Message: "Register completed", StatusCode: 200}, map[string]interface{}{"function": "Create"})
}

func (uc *HttpUser) Get(c *fiber.Ctx) error {
	return nil
}

func (uc *HttpUser) GetAll(c *fiber.Ctx) error {
	return nil
}

func (uc *HttpUser) Update(c *fiber.Ctx) error {
	return nil
}

func (uc *HttpUser) Delete(c *fiber.Ctx) error {
	return nil
}
