package repository

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"

	"github.com/RJHsiao/api-exercise-go/database"
	"github.com/RJHsiao/api-exercise-go/utilities"
)

// AssignSessionKey assign session key when user login
func AssignSessionKey(userID primitive.ObjectID) (string, error) {
	sessionKey := utilities.GetSha256SumFromString(userID.Hex() + time.Now().String())
	session := database.Session{
		ID:          primitive.NewObjectID(),
		SessionKey:  sessionKey,
		UserID:      userID,
		ExpiredTime: primitive.DateTime(time.Now().AddDate(0, 0, 7).UnixNano()),
	}
	_, err := database.CollectionSessions.InsertOne(context.TODO(), session)
	return sessionKey, err
}

// RevokeSessionKey revoke session key when user logout
func RevokeSessionKey(sessionKey string) {
	filter := bson.D{{Key: "sessionKey", Value: sessionKey}}
	database.CollectionSessions.FindOneAndDelete(context.TODO(), filter)
}

// FindSessionByKey get session info by session key
func FindSessionByKey(sessionKey string) (database.Session, error) {
	filter := bson.D{{Key: "sessionKey", Value: sessionKey}}
	var session database.Session
	err := database.CollectionSessions.FindOne(context.TODO(), filter).Decode(&session)
	return session, err
}
