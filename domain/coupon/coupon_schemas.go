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
	ID           primitive.ObjectID `bson:"_id"`
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

// NewCouponRequest is a new coupon DTO
type NewCouponRequest struct {
	ID          primitive.ObjectID         `json:"id" `
	Description string                     `json:"description" validate:"required"`
	Amount      float32                    `json:"amount" validate:"min=0"`
	Percentage  uint8                      `json:"percentage" validate:"min=0"`
	Constraint  NewCouponConstraintRequest `json:"constraint"`
	CouponType  string                     `json:"coupon_type"`
}

// NewCouponResponse is a new coupon DTO
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

// NewCouponConstraintRequest is a new coupon constraint DTO
type NewCouponConstraintRequest struct {
	ValidityFrom time.Time `json:"validity_from" validate:"required"`
	ValidityTo   time.Time `json:"validity_to" validate:"required"`
	TotalUsage   int32     `json:"total_usage" validate:"min=1,required"`
	MaxUsage     int32     `json:"max_usage" validate:"min=1,required"`
	MaxAmount    float32   `json:"max_amount" validate:"min=1,required"`
	MinItems     int32     `json:"min_items" validate:"min=1,required"`
	MaxItems     int32     `json:"max_items" validate:"min=1,required"`
	Combinable   bool      `json:"combinable" validate:"required"`
}

// NewCouponConstraintResponse is a new coupon constraint DTO
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
