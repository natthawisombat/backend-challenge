package usecases

import (
	"backend-challenge/entities"
	"context"
)

type userRepository interface {
	Register(user entities.User, ctx context.Context) error
	CheckDuplicateUser(email string, ctx context.Context) error
	Login(login entities.Login, ctx context.Context) (string, error)
	GetUserAll(ctx context.Context) ([]entities.User, error)
	GetUser(userId string, ctx context.Context) (result entities.User, err error)
	DeleteUser(userId string, ctx context.Context) error
	UpdateUser(userId string, data entities.UpdateUserRequest, ctx context.Context) error
}
