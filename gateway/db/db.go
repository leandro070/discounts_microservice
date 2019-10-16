package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var dbInstance *mongo.Client
var dbContext context.Context

// GetMongo return MongoDB instance
func GetMongo() (*mongo.Client, context.Context) {
	return dbInstance, dbContext
}

// InitMongoClient inicialize MongoDB
func InitMongoClient() error {
	clientOpt := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.NewClient(clientOpt)
	if err != nil {
		return err
	}

	err = client.Connect(context.Background())
	if err != nil {
		return err
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return err
	}
	dbInstance = client
	return nil
}
