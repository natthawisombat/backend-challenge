package usecases

import (
	"backend-challenge/entities"
	"context"
)

type userRepository interface {
	Register(user entities.User, ctx context.Context) error
	CheckDuplicateUser(email string, ctx context.Context) error
	Login(login entities.Login, ctx context.Context) (string, error)
}
