package database

import (
	"context"
	"log"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/x/network/connstring"

	"github.com/RJHsiao/api-exercise-go/config"
)

var (
	client *mongo.Client
	name   string

	// CollectionUsers a collection to store users' info
	CollectionUsers *mongo.Collection
	// CollectionSessions a collection to store login session
	CollectionSessions *mongo.Collection
)

// Connect database.
func Connect() error {
	var cs connstring.ConnString
	var c *mongo.Client
	var err error

	cs, err = connstring.Parse(config.GetConfig().Database)
	if err != nil {
		return err
	}
	name = cs.Database
	log.Println("Connect database:", cs)

	c, err = mongo.Connect(context.TODO(), config.GetConfig().Database)
	if err != nil {
		return err
	}
	client = c
	// Check connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}

	CollectionUsers = client.Database(name).Collection(collectionNameUsers)
	CollectionSessions = client.Database(name).Collection(collectionNameSessions)

	return nil
}

// Disconnect database
func Disconnect() {
	err := client.Disconnect(context.TODO())
	log.Println("Disconnect database with error:", err)
}
