package coupon

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CouponSchema is a struct to coupon model
type CouponSchema struct {
	ID           primitive.ObjectID `bson:"_id"`
	Description  string             `bson:"description" validate:"required"`
	Code         string             `bson:"code"`
	Amount       float32            `bson:"amount"`
	Percentage   uint8              `bson:"percentage"`
	IsEnable     bool               `bson:"is_enable"`
	ConstraintID primitive.ObjectID `bson:"constraint_id"`
	CouponType   string             `bson:"coupon_type"`
	UpdatedAt    time.Time          `bson:"update_at"`
	CreatedAt    time.Time          `bson:"created_at"`
}

// CouponContraint is a coupon limitation
type CouponContraint struct {
	ID           primitive.ObjectID `bson:"_id" `
	ValidityFrom time.Time          `bson:"validity_from"`
	ValidityTo   time.Time          `bson:"validity_to"`
	TotalUsage   int32              `bson:"total_usage"`
	MaxUsage     int32              `bson:"max_usage"`
	MaxAmount    float32            `bson:"max_amount"`
	MinItems     int32              `bson:"min_items"`
	MaxItems     int32              `bson:"max_items"`
	Combinable   bool               `bson:"combinable"`
	UpdatedAt    time.Time          `bson:"update_at"`
	CreatedAt    time.Time          `bson:"created_at"`
}

// NewCouponRequest
type NewCouponRequest struct {
	ID          primitive.ObjectID         `json:"id" `
	Description string                     `json:"description"`
	Amount      float32                    `json:"amount"`
	Percentage  uint8                      `json:"percentage"`
	Constraint  NewCouponConstraintRequest `json:"constraint"`
	CouponType  string                     `json:"coupon_type"`
}

// NewCouponResponse
type NewCouponResponse struct {
	ID          primitive.ObjectID          `json:"id"`
	Description string                      `json:"description"`
	Amount      float32                     `json:"amount"`
	Percentage  uint8                       `json:"percentage"`
	IsEnable    bool                        `json:"is_enable"`
	Code        string                      `json:"code"`
	Constraint  NewCouponConstraintResponse `json:"constraint"`
	CouponType  string                      `json:"coupon_type"`
	UpdatedAt   time.Time                   `json:"update_at"`
	CreatedAt   time.Time                   `json:"created_at"`
}

// NewCouponConstraintRequest
type NewCouponConstraintRequest struct {
	ValidityFrom time.Time `json:"validity_from"`
	ValidityTo   time.Time `json:"validity_to"`
	TotalUsage   int32     `json:"total_usage"`
	MaxUsage     int32     `json:"max_usage"`
	MaxAmount    float32   `json:"max_amount"`
	MinItems     int32     `json:"min_items"`
	MaxItems     int32     `json:"max_items"`
	Combinable   bool      `json:"combinable"`
}

// NewCouponConstraintResponse
type NewCouponConstraintResponse struct {
	ID           primitive.ObjectID `json:"id"`
	ValidityFrom time.Time          `json:"validity_from"`
	ValidityTo   time.Time          `json:"validity_to"`
	TotalUsage   int32              `json:"total_usage"`
	MaxUsage     int32              `json:"max_usage"`
	MaxAmount    float32            `json:"max_amount"`
	MinItems     int32              `json:"min_items"`
	MaxItems     int32              `json:"max_items"`
	Combinable   bool               `json:"combinable"`
}
