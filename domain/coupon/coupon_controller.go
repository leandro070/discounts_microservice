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
		"result": res,
	})
}

// GetCoupon se encargará de recibir un código de descuento, validar la existencia y vigencia del cupón.
func GetCoupon(c *gin.Context) {

	couponID := c.Param("id")

	couponService, err := NewService()
	if err != nil {
		errors.Handle(c, err)
		return
	}

	res, err := couponService.GetCoupon(couponID)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.JSON(200, gin.H{
		"result": res,
	})
}

// AnnulCoupon se encargará de recibir un código de descuento y darlo de baja.
func AnnulCoupon(c *gin.Context) {

	couponID := c.Param("id")

	couponService, err := NewService()
	if err != nil {
		errors.Handle(c, err)
		return
	}

	err = couponService.AnnulCoupon(couponID)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.JSON(200, gin.H{
		"success": true,
	})
}
