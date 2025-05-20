package adapters

import (
	"backend-challenge/entities"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	db *mongo.Database
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{db: db}
}

func (rp *MongoRepository) Register(user entities.User, ctx context.Context) error {
	coll := rp.db.Collection("user")
	_, err := coll.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (rp *MongoRepository) CheckDuplicateUser(email string, ctx context.Context) error {
	coll := rp.db.Collection("user")
	filter := bson.M{
		"email": email,
	}

	var results entities.User
	err := coll.FindOne(ctx, filter).Decode(&results)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
		return err
	}

	return fmt.Errorf("email already exists")
}

func (rp *MongoRepository) Login(login entities.Login, ctx context.Context) (string, error) {
	coll := rp.db.Collection("user")
	filter := bson.M{
		"$and": []bson.M{
			{"email": login.Email},
			{"password": login.Password},
		},
	}
	var results entities.User
	err := coll.FindOne(ctx, filter).Decode(&results)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", fmt.Errorf("Email or Password was wrong.")
		}
		return "", err
	}
	return results.ID.Hex(), nil
}
