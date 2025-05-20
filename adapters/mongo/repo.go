package adapters

import (
	"backend-challenge/entities"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (rp *MongoRepository) GetUserAll(ctx context.Context) ([]entities.User, error) {
	coll := rp.db.Collection("user")
	var result []entities.User
	cursor, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (rp *MongoRepository) GetUser(userId string, ctx context.Context) (result entities.User, err error) {
	coll := rp.db.Collection("user")
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return result, fmt.Errorf("invalid user ID format: %w", err)
	}

	filter := bson.M{
		"_id": oid,
	}
	if err := coll.FindOne(ctx, filter).Decode(&result); err != nil {
		return result, err
	}

	return result, err
}

func (rp *MongoRepository) DeleteUser(userId string, ctx context.Context) error {
	coll := rp.db.Collection("user")
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	filter := bson.M{
		"_id": oid,
	}

	if _, err := coll.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}

func (rp *MongoRepository) UpdateUser(userId string, data entities.UpdateUserRequest, ctx context.Context) error {
	coll := rp.db.Collection("user")
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	filter := bson.M{
		"_id": oid,
	}

	update := bson.M{}
	if data.Email != "" {
		update["email"] = data.Email
	}

	if data.Name != "" {
		update["name"] = data.Name
	}

	_, err = coll.UpdateOne(ctx, filter, bson.M{"$set": update})
	return err
}
