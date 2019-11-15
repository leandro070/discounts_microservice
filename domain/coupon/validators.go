package coupon

import (
	"fmt"
	"time"

	"github.com/leandro070/discounts_microservice/utils/errors"

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

func (c CouponSchema) validate() error {

	if c.IsEnable == false {
		return errors.NewValidationField("is_enable", "cupon disabled")
	}

	if (c.Amount == nil || c.Amount == 0) && c.CouponType == "fixed_amount" {
		return errors.NewValidationField("amount", "Coupon type 'fixed_amount' must have amount higher than zero")
	}

	if (c.Percentage == nil || c.Percentage == 0) && c.CouponType == "percentage" {
		return errors.NewValidationField("percentage", "Coupon type 'percentage' must have percentage higher than zero")
	}

	return nil
}

func (c CouponContraint) validate(items int) error {

	if c.IsEnable == false {
		return errors.NewValidationField("is_enable", "constraint disabled")
	}
	// in range time
	datetime := time.Now()
	if c.ValidityFrom.After(datetime) {
		return errors.NewValidationField("validity_from", "Coupon invalid")
	}
	if c.ValidityTo.Before(datetime) {
		return errors.NewValidationField("validity_to", "Coupon invalid")
	}

	// Max use
	if c.MaxUsage != 0 && c.TotalUsage == c.MaxUsage {
		return errors.NewValidationField("max_usage", "Max usage for coupon")
	}

	if c.MinItems > int32(items) {
		return errors.NewValidationField("min_items", fmt.Sprintf("Coupon need %d items for apply", c.MinItems))
	}

	return nil
}

func (c CouponContraint) validateNew() error {

	if c.ValidityFrom.After(c.ValidityTo) {
		return errors.NewValidationField("validity_from", "Invalid dates")
	}

	// in range time
	datetime := time.Now()
	if c.ValidityFrom.Before(datetime) {
		return errors.NewValidationField("validity_from", "Date before today")
	}
	if c.ValidityTo.Before(datetime) {
		return errors.NewValidationField("validity_to", "Date before today")
	}

	// Max use
	if c.MaxUsage < 0 {
		return errors.NewValidationField("max_usage", "Max usage must be positive number or zero")
	}

	if c.MaxItems < 0 {
		return errors.NewValidationField("max_items", "Max items must be positive number or zero")
	}

	if c.MinItems < 0 {
		return errors.NewValidationField("min_items", "Min items must be positive number or zero")
	}

	return nil
}
