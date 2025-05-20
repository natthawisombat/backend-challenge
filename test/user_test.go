package user_test

import (
	"backend-challenge/entities"
	"backend-challenge/usecases"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) Register(user entities.User, ctx context.Context) error {
	args := m.Called(user, ctx)
	return args.Error(0)
}
func (m *mockUserRepo) CheckDuplicateUser(email string, ctx context.Context) error {
	args := m.Called(email, ctx)
	return args.Error(0)
}
func (m *mockUserRepo) Login(login entities.Login, ctx context.Context) (string, error) {
	args := m.Called(login, ctx)
	return args.String(0), args.Error(1)
}
func (m *mockUserRepo) GetUserAll(ctx context.Context) ([]entities.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entities.User), args.Error(1)
}
func (m *mockUserRepo) GetUser(id string, ctx context.Context) (entities.User, error) {
	args := m.Called(id, ctx)
	return args.Get(0).(entities.User), args.Error(1)
}
func (m *mockUserRepo) DeleteUser(id string, ctx context.Context) error {
	args := m.Called(id, ctx)
	return args.Error(0)
}
func (m *mockUserRepo) UpdateUser(id string, input entities.UpdateUserRequest, ctx context.Context) error {
	args := m.Called(id, input, ctx)
	return args.Error(0)
}

func setupTestApp(handler usecases.HttpUser) *fiber.App {
	app := fiber.New()
	app.Post("/auth/register", handler.Create)
	app.Post("/auth/login", handler.Login)
	app.Get("/users", handler.GetAll)
	app.Get("/users/:id", handler.Get)
	app.Patch("/users/:id", handler.Update)
	app.Delete("/users/:id", handler.Delete)
	return app
}

func TestRegisterUser(t *testing.T) {
	repo := new(mockUserRepo)
	h := usecases.NewHttpUser(validator.New(), repo)
	app := setupTestApp(h)

	body := `{"name":"Tee","email":"tee@email.com","password":"123456"}`
	repo.On("CheckDuplicateUser", "tee@email.com", mock.Anything).Return(nil)
	repo.On("Register", mock.AnythingOfType("entities.User"), mock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	repo.AssertExpectations(t)
}

func TestRegisterUser_Duplicate(t *testing.T) {
	repo := new(mockUserRepo)
	h := usecases.NewHttpUser(validator.New(), repo)
	app := setupTestApp(h)

	body := `{"name":"Tee","email":"tee@email.com","password":"123456"}`
	repo.On("CheckDuplicateUser", "tee@email.com", mock.Anything).Return(errors.New("duplicate"))

	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestLoginSuccess(t *testing.T) {
	repo := new(mockUserRepo)
	h := usecases.NewHttpUser(validator.New(), repo)
	app := setupTestApp(h)

	input := `{"email":"a@b.com","password":"123456"}`
	repo.On("Login", mock.AnythingOfType("entities.Login"), mock.Anything).Return("userid123", nil)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(input))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestGetAllUsers(t *testing.T) {
	repo := new(mockUserRepo)
	h := usecases.NewHttpUser(validator.New(), repo)
	app := setupTestApp(h)

	repo.On("GetUserAll", mock.Anything).Return([]entities.User{{Name: "Tee", Email: "a@b.com"}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestGetUser(t *testing.T) {
	repo := new(mockUserRepo)
	h := usecases.NewHttpUser(validator.New(), repo)
	app := setupTestApp(h)

	repo.On("GetUser", "abc123", mock.Anything).Return(entities.User{Name: "Tee"}, nil)

	req := httptest.NewRequest(http.MethodGet, "/users/abc123", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestUpdateUser(t *testing.T) {
	repo := new(mockUserRepo)
	h := usecases.NewHttpUser(validator.New(), repo)
	app := setupTestApp(h)

	repo.On("CheckDuplicateUser", "a@b.com", mock.Anything).Return(nil)
	repo.On("UpdateUser", "abc123", mock.AnythingOfType("entities.UpdateUserRequest"), mock.Anything).Return(nil)

	body := `{"name":"Tee","email":"a@b.com"}`
	req := httptest.NewRequest(http.MethodPatch, "/users/abc123", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestDeleteUser(t *testing.T) {
	repo := new(mockUserRepo)
	h := usecases.NewHttpUser(validator.New(), repo)
	app := setupTestApp(h)

	repo.On("DeleteUser", "abc123", mock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/users/abc123", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)
}
