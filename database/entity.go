package database

import (
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

const (
	collectionNameUsers    = "users"
	collectionNameSessions = "sessions"
)

// User user entity
type User struct {
	ID         primitive.ObjectID `bson:"_id"`
	Name       string             `bson:"name"`
	Email      string             `bson:"email"`
	Password   string             `bson:"password"`
	UpdateTime primitive.DateTime `bson:"updateAt"`
}

// Session session entity
type Session struct {
	ID          primitive.ObjectID `bson:"_id"`
	SessionKey  string             `bson:"sessionKey"`
	UserID      primitive.ObjectID `bson:"userId"`
	ExpiredTime primitive.DateTime `bson:"expiredAt"`
}
