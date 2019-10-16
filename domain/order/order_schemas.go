package order

import (
	"time"

	"github.com/leandro070/discounts_microservice/domain/coupon"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID        primitive.ObjectID `json: "id" bson: "_id"`
	Status    string             `json: "status" bson: "status"`
	UserId    primitive.ObjectID `json: "user_id" bson: "user_id"`
	CartId    primitive.ObjectID `json: "cart_id" bson: "cart_id"`
	UpdatedAt time.Time          `json: "updated" bson: "update_at"`
	CreatedAt time.Time          `json: "created" bson: "created_at`
	Coupons   []coupon.Coupon    `json: "cuopons" bson: "coupons"`
}
