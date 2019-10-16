package coupon

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/leandro070/discounts_microservice/gateway/db"
	"go.mongodb.org/mongo-driver/bson"
)

// NewCoupon se encargará de crear un cupón y sus restricciones asociadas
func NewCoupon(c *gin.Context) {
	s, ctx := db.GetMongo()

	col := s.Database("discount").Collection("coupons")

	filter := bson.M{"coupons": bson.M{"$elemMatch": bson.M{"$eq": "golang"}}}

	// find one document
	var coupon Coupon
	if err := col.FindOne(ctx, filter).Decode(&coupon); err != nil {
		log.Panic(err)
	}
	fmt.Printf("coupon: %+v\n", c)

	c.JSON(200, gin.H{
		"coupons": "go",
	})
}
