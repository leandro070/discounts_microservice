package coupon

import (
	"fmt"

	"github.com/leandro070/discounts_microservice/utils/errors"
)

func validate(obj interface{}) error {

	switch obj.(type) {
	case Coupon:
		return validateNewCoupon(obj)
	}

	return nil
}

func validateNewCoupon(obj interface{}) error {
	coupon, ok := obj.(Coupon)
	if ok != true {
		return errors.NewCustom(400, "Coupon binding invalid")
	}

	fmt.Println(coupon.Description)
	return nil
}
