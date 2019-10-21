package coupon

import (
	"github.com/gin-gonic/gin"
	"github.com/leandro070/discounts_microservice/utils/errors"
)

// NewCoupon se encargará de crear un cupón y sus restricciones asociadas
func NewCoupon(c *gin.Context) {

	var coupon NewCouponRequest
	err := c.ShouldBindJSON(&coupon)
	if err != nil {
		errors.Handle(c, err)
		panic(err.Error)
		return
	}

	couponService, err := NewService()
	if err != nil {
		errors.Handle(c, err)
		return
	}

	res, err := couponService.NewCoupon(&coupon)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.JSON(200, gin.H{
		"coupons": res,
	})

}
