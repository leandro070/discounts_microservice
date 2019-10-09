package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cMongo := getMongoClient()

	r := gin.Default()

	v1 := r.Group("/v1")
	{
		v1.GET("/", Root)
	}

	r.Run(":3030")
}

func getMongoClient() *mongo.Client {
	clientOpt := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't connect to the Mongo database", err)
	} else {
		log.Println("MongoDb connected!")
	}

	return client
}
