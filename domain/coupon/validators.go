package coupon

import (
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
