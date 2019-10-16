package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/leandro070/discounts_microservice/domain/coupon"
	"github.com/leandro070/discounts_microservice/gateway/rabbit"
	"github.com/leandro070/discounts_microservice/utils/env"
	"github.com/leandro070/discounts_microservice/utils/errors"
	"github.com/leandro070/discounts_microservice/utils/security"
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
	v1.Use(validateAuthentication())
	{
		v1.Use(static.Serve("/", static.LocalFile(env.Get().WWWWPath, true)))
		v1.POST("/coupons", coupon.NewCoupon)
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

// get token from Authorization header
func getTokenHeader(c *gin.Context) (string, error) {
	tokenString := c.GetHeader("Authorization")
	if strings.Index(tokenString, "bearer ") != 0 {
		return "", errors.Unauthorized
	}
	return tokenString[7:], nil
}

func validateAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := getTokenHeader(c)
		if err != nil {
			errors.Handle(c, errors.Unauthorized)
			c.Abort()
			return
		}

		if _, err = security.Validate(tokenString); err != nil {
			errors.Handle(c, errors.Unauthorized)
			c.Abort()
			return
		}
	}
}
