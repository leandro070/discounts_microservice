package coupon

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Coupon is a struct to coupon model
type Coupon struct {
	ID          primitive.ObjectID `json: "id" bson: "_id"`
	Description string             `json: "description" bson:"description"`
	Code        string             `json: "code" bson: "code"`
	Amount      float32            `json: "amount" bson: "amount"`
	Percentage  uint8              `json: "percentage" bson: "percentage"`
	IsEnable    bool               `json: "is_enable" bson: "is_enable"`
	Constraint  Contraint          `json: "constraint"`
	CouponType  CouponType         `json: "coupon_type"`
	UpdatedAt   time.Time          `json: "updated" bson: "update_at"`
	CreatedAt   time.Time          `json: "created" bson: "created_at`
}

type CouponType struct {
	ID   primitive.ObjectID `json: "id" bson: "_id"`
	Name string             `json: "name" bson: "name"`
}

// Contraint is a coupon limitation
type Contraint struct {
	ID           primitive.ObjectID `json: "id" bson: "_id"`
	ValidityFrom time.Time          `json: "validity_from" bson: "validity_from"`
	ValidityTo   time.Time          `json: "validity_to" bson: "validity_to"`
	TotalUsage   int32              `json: "total_usage" bson: "total_usage"`
	MaxUsage     int32              `json: "max_usage" bson: "max_usage"`
	MaxAmount    float32            `json: "max_amount" bson: "max_amount"`
	MinItems     int32              `json: "min_items" bson: "min_items"`
	MaxItems     int32              `json: "max_items" bson: "max_items"`
	Combinable   bool               `json: "combinable" bson: "combinable"`
	UpdatedAt    time.Time          `json: "updated" bson: "update_at"`
	CreatedAt    time.Time          `json: "created" bson: "created_at`
}
