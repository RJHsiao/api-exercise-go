package repository

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"

	"github.com/RJHsiao/api-exercise-go/database"
	"github.com/RJHsiao/api-exercise-go/models"
	"github.com/RJHsiao/api-exercise-go/utilities"
)

// IsEmailRegistered check given email is registered or not
func IsEmailRegistered(email string) (bool, error) {
	filter := bson.D{{Key: "email", Value: email}}
	n, err := database.CollectionUsers.Count(context.TODO(), filter)
	if err != nil {
		return false, err
	}
	return (n > 0), nil
}

// AddUser add user
func AddUser(userForm models.RequestEditUserForm) error {
	newUser := database.User{
		ID:         primitive.NewObjectID(),
		Name:       userForm.Name,
		Email:      userForm.Email,
		Password:   utilities.GetSha256SumFromString(userForm.Password),
		UpdateTime: primitive.DateTime(time.Now().Unix() * 1000),
	}

	_, err := database.CollectionUsers.InsertOne(context.TODO(), newUser)
	return err
}

// UpdateUser edit user
func UpdateUser(newUser database.User) error {
	filter := bson.D{{Key: "_id", Value: newUser.ID}}
	_, err := database.CollectionUsers.ReplaceOne(context.TODO(), filter, newUser)
	return err
}

// FindUserByLoginForm Check user exist and password matched
func FindUserByLoginForm(form models.RequestLoginForm) (*database.User, error) {
	filter := bson.D{
		{Key: "email", Value: form.Email},
		{Key: "password", Value: utilities.GetSha256SumFromString(form.Password)},
	}
	var user database.User
	err := database.CollectionUsers.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindUserByObjectID get user info by user's objectId
func FindUserByObjectID(objectID primitive.ObjectID) (*database.User, error) {
	filter := bson.D{{Key: "_id", Value: objectID}}
	var user database.User
	err := database.CollectionUsers.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
