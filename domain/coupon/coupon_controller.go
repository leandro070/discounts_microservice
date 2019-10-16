package coupon

import (
	"github.com/gin-gonic/gin"
)

func NewCoupon(c *gin.Context) {
	c.JSON(200, gin.H{
		"ok": "ok",
	})
}
