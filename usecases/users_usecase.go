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

func (uc *HttpUser) Login(c *fiber.Ctx) error {
	var bodyRequest entities.Login
	if err := c.BodyParser(&bodyRequest); err != nil {
		return handlers.Response(c, entities.Response{Status: "ER", ErrorMessage: err.Error(), ErrorCode: "ER400", StatusCode: 400}, map[string]interface{}{"function": "Login"})
	}

	// validate request body
	err := uc.validate.Struct(bodyRequest)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return handlers.Response(c, entities.Response{Status: "ER", ErrorMessage: err.Error(), ErrorCode: "ER400", StatusCode: 400}, map[string]interface{}{"function": "Login"})
		}
	}

	//hash
	bodyRequest.Password = utils.Hash(bodyRequest.Password)
	userId, err := uc.repo.Login(bodyRequest, c.UserContext())
	if err != nil {
		return handlers.Response(c, entities.Response{Status: "ER", ErrorMessage: err.Error(), ErrorCode: "ER400", StatusCode: 400}, map[string]interface{}{"function": "Login"})
	}

	authKey, err := utils.GenerateToken(userId, time.Duration(24*time.Hour))
	if err != nil {
		return handlers.Response(c, entities.Response{Status: "ER", ErrorMessage: err.Error(), ErrorCode: "ER400", StatusCode: 400}, map[string]interface{}{"function": "Login"})
	}

	return handlers.Response(c,
		entities.Response{Status: "OK", Message: "Success", StatusCode: 200, Data: map[string]interface{}{
			"accessToken": authKey,
		}}, map[string]interface{}{"function": "Login"})
}

func (uc *HttpUser) Create(c *fiber.Ctx) error {
	var bodyRequest entities.User
	if err := c.BodyParser(&bodyRequest); err != nil {
		return handlers.Response(c, entities.Response{Status: "ER", ErrorMessage: err.Error(), ErrorCode: "ER400", StatusCode: 400}, map[string]interface{}{"function": "Create"})
	}

	// validate request body
	err := uc.validate.Struct(bodyRequest)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return handlers.Response(c, entities.Response{Status: "ER", ErrorMessage: err.Error(), ErrorCode: "ER400", StatusCode: 400}, map[string]interface{}{"function": "Create"})
		}
	}

	// check duplicate
	if err := uc.repo.CheckDuplicateUser(bodyRequest.Email, c.UserContext()); err != nil {
		return handlers.Response(c, entities.Response{Status: "ER", ErrorMessage: err.Error(), ErrorCode: "ER400", StatusCode: 400}, map[string]interface{}{"function": "Create"})
	}
	// hash password
	bodyRequest.Password = utils.Hash(bodyRequest.Password)
	bodyRequest.CreatedAt = time.Now().Add(7 * time.Hour)

	if err := uc.repo.Register(bodyRequest, c.UserContext()); err != nil {
		return handlers.Response(c, entities.Response{Status: "ER", ErrorMessage: err.Error(), ErrorCode: "ER400", StatusCode: 400}, map[string]interface{}{"function": "Create"})
	}

	return handlers.Response(c, entities.Response{Status: "OK", Message: "Register completed", StatusCode: 200}, map[string]interface{}{"function": "Create"})
}

func (uc *HttpUser) Get(c *fiber.Ctx) error {
	userId := c.Params("id")
	user, err := uc.repo.GetUser(userId, c.UserContext())
	if err != nil {
		return handlers.Response(c, entities.Response{Status: "ER", ErrorMessage: err.Error(), ErrorCode: "ER404", StatusCode: 200}, map[string]interface{}{"function": "Get"})
	}

	return handlers.Response(c, entities.Response{Status: "OK", Data: user, StatusCode: 200}, map[string]interface{}{"function": "Get"})
}

func (uc *HttpUser) GetAll(c *fiber.Ctx) error {
	// check duplicate
	users, err := uc.repo.GetUserAll(c.UserContext())
	if err != nil {
		return handlers.Response(c, entities.Response{Status: "ER", ErrorMessage: err.Error(), ErrorCode: "ER400", StatusCode: 400}, map[string]interface{}{"function": "GetAll"})
	}

	return handlers.Response(c, entities.Response{Status: "OK", Data: users, StatusCode: 200}, map[string]interface{}{"function": "GetAll"})
}

func (uc *HttpUser) Update(c *fiber.Ctx) error {
	userId := c.Params("id")
	var bodyRequest entities.UpdateUserRequest
	if err := c.BodyParser(&bodyRequest); err != nil {
		return handlers.Response(c, entities.Response{Status: "ER", ErrorMessage: err.Error(), ErrorCode: "ER400", StatusCode: 400}, map[string]interface{}{"function": "Update"})
	}

	// validate request body
	err := uc.validate.Struct(bodyRequest)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return handlers.Response(c, entities.Response{Status: "ER", ErrorMessage: err.Error(), ErrorCode: "ER400", StatusCode: 400}, map[string]interface{}{"function": "Update"})
		}
	}

	if err := uc.repo.CheckDuplicateUser(bodyRequest.Email, c.UserContext()); err != nil {
		return handlers.Response(c, entities.Response{Status: "ER", ErrorMessage: err.Error(), ErrorCode: "ER400", StatusCode: 400}, map[string]interface{}{"function": "Update"})
	}

	if err := uc.repo.UpdateUser(userId, bodyRequest, c.UserContext()); err != nil {
		return handlers.Response(c, entities.Response{Status: "ER", ErrorMessage: err.Error(), ErrorCode: "ER400", StatusCode: 400}, map[string]interface{}{"function": "Update"})
	}

	return handlers.Response(c, entities.Response{Status: "OK", Message: "Update success", StatusCode: 200}, map[string]interface{}{"function": "Update"})
}

func (uc *HttpUser) Delete(c *fiber.Ctx) error {
	userId := c.Params("id")
	if err := uc.repo.DeleteUser(userId, c.UserContext()); err != nil {
		return handlers.Response(c, entities.Response{Status: "ER", ErrorMessage: err.Error(), ErrorCode: "ER500", StatusCode: 500}, map[string]interface{}{"function": "Delete"})
	}
	return handlers.Response(c, entities.Response{Status: "OK", Message: "Delete success", StatusCode: 200}, map[string]interface{}{"function": "Delete"})
}
