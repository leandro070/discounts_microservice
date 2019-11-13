package coupon

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func validateCode(s *ServiceCupon, code string) bool {
	coupon, err := s.repo.FindByCodeCoupon(code)
	if err != nil {
		return true // responde por error cuando no hay documento con ese code
	}
	if coupon.ID != primitive.NilObjectID {
		return false
	}
	return false
}

func (c CouponContraint) validate(items int) error {

	// in range time
	datetime := time.Now()
	if c.ValidityFrom.After(datetime) || c.ValidityTo.Before(datetime) {
		return fmt.Errorf("Not valid coupon")
	}
	// Max use
	if c.MaxUsage != 0 && c.TotalUsage == c.MaxUsage {
		return fmt.Errorf("Max usage for coupon")
	}

	if c.MinItems > int32(items) {
		return fmt.Errorf("Coupon need %d items for apply", c.MinItems)
	}

	return nil
}

func (c CouponContraint) validateNew() error {

	if c.ValidityFrom.After(c.ValidityTo) {
		return fmt.Errorf("Invalid dates")
	}

	// in range time
	datetime := time.Now()
	if c.ValidityFrom.Before(datetime) || c.ValidityTo.Before(datetime) {
		return fmt.Errorf("Invalid dates")
	}

	// Max use
	if c.MaxUsage < 0 {
		return fmt.Errorf("Max usage must be positive number or zero")
	}

	if c.MaxItems < 0 {
		return fmt.Errorf("Max items must be positive number or zero")
	}

	if c.MinItems < 0 {
		return fmt.Errorf("Min items must be positive number or zero")
	}

	return nil
}
