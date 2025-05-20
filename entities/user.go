package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name" validate:"required"`
	Email     string             `bson:"email" json:"email" validate:"required,email"`
	Password  string             `bson:"password" json:"password" validate:"required"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}

type Login struct {
	Email    string `bson:"email" json:"email" validate:"required,email"`
	Password string `bson:"password" json:"password" validate:"required"`
}
