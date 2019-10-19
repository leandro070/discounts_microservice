package coupon

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/leandro070/discounts_microservice/gateway/db"
	"github.com/leandro070/discounts_microservice/utils/errors"
)

// NewCoupon se encargará de crear un cupón y sus restricciones asociadas
func NewCoupon(c *gin.Context) {
	s, ctx := db.GetMongo()

	var coupon Coupon
	err := c.BindJSON(&coupon)
	err = validate(coupon)
	if err != nil {
		errors.Handle(c, err)
	}

	col := s.Database("discount").Collection("coupons")

	result, err := col.InsertOne(ctx, coupon)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("coupon: %+v\n", coupon)

	c.JSON(200, gin.H{
		"coupons": result,
	})
}
