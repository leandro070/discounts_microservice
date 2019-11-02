package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/leandro070/discounts_microservice/domain/coupon"
	"github.com/leandro070/discounts_microservice/gateway/db"
	"github.com/leandro070/discounts_microservice/utils/env"
	"github.com/leandro070/discounts_microservice/utils/errors"
	"github.com/leandro070/discounts_microservice/utils/security"
)

func main() {
	if len(os.Args) > 1 {
		env.Load(os.Args[1])
	}

	coupon.RabbitInit()

	err := db.InitMongoClient()
	if err != nil {
		log.Fatal("Couldn't connect to the Mongo database", err)
	}
	log.Println("MongoDb connected!")

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "GET", "POST"},
		AllowHeaders:     []string{"Origin, Authorization, Content-Type"},
		ExposeHeaders:    []string{""},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(static.Serve("/", static.LocalFile(env.Get().WWWWPath, true)))
	v1 := r.Group("/v1")
	{
		v1.POST("/coupons", coupon.NewCoupon)
		v1.GET("/coupons/:id", coupon.GetCoupon)
		v1.DELETE("/coupons/:id", coupon.AnnulCoupon)
	}

	r.Run(":3030")
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
