package service

import (
	"context"

	"github.com/MusaSSH/go-mongo-rest-example/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	collection *mongo.Collection
}

func (s UserService) Insert(document models.User) (primitive.ObjectID, error) {
	result, err := s.collection.InsertOne(context.TODO(), document)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (s UserService) Get(filter models.User) (models.User, error) {
	result := s.collection.FindOne(context.TODO(), filter)

	if result.Err() != nil {
		return models.User{}, result.Err()
	}

	var user models.User
	err := result.Decode(&user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s UserService) Update(objectid primitive.ObjectID, update models.User) error {
	_, err := s.collection.UpdateByID(context.TODO(), objectid, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		return err
	}
	return nil
}

func NewUserService(client *mongo.Client) UserService {
	return UserService{
		collection: client.Database("mydb").Collection("user"),
	}
}
