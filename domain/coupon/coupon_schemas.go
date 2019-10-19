package coupon

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Coupon is a struct to coupon model
type Coupon struct {
	ID          primitive.ObjectID `bson: "_id"`
	Description string             `bson: "description" validate: "required"`
	Code        string             `bson: "code"`
	Amount      float32            `bson: "amount"`
	Percentage  uint8              `bson: "percentage"`
	IsEnable    bool               `bson: "is_enable"`
	Constraint  CouponContraint    `bson: "constraint"`
	CouponType  CouponType         `bson: "coupon_type"`
	UpdatedAt   time.Time          `bson: "update_at"`
	CreatedAt   time.Time          `bson: "created_at`
}

type CouponType struct {
	ID   primitive.ObjectID `bson: "_id"`
	Name string             `bson: "name"`
}

// CouponContraint is a coupon limitation
type CouponContraint struct {
	ID           primitive.ObjectID `bson: "_id"`
	ValidityFrom time.Time          `bson: "validity_from"`
	ValidityTo   time.Time          `bson: "validity_to"`
	TotalUsage   int32              `bson: "total_usage"`
	MaxUsage     int32              `bson: "max_usage"`
	MaxAmount    float32            `bson: "max_amount"`
	MinItems     int32              `bson: "min_items"`
	MaxItems     int32              `bson: "max_items"`
	Combinable   bool               `bson: "combinable"`
	UpdatedAt    time.Time          `bson: "update_at"`
	CreatedAt    time.Time          `bson: "created_at`
}
