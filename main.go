package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/leandro070/discounts_microservice/utils/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	if len(os.Args) > 1 {
		env.Load(os.Args[1])
	}

	rabbit.Init()
	
	client := getMongoClient()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "GET", "POST"},
		AllowHeaders:     []string{"Origin, Authorization, Content-Type"},
		ExposeHeaders:    []string{""},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	
	v1 := r.Group("/v1")
	{
		v1.Use(static.Serve("/", static.LocalFile(env.Get().WWWWPath, true)))
		
	}

	r.Run(":3030")
}

func getMongoClient() *mongo.Client {
	clientOpt := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.NewClient(clientOpt)
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
