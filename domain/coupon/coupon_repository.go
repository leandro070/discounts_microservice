package coupon

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/leandro070/discounts_microservice/gateway/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	InsertCoupon(coupon CouponSchema) (CouponSchema, error)
	UpdateCoupon(coupon CouponSchema) (CouponSchema, error)
	FindAllCoupons() ([]CouponSchema, error)
	FindByIDCoupon(couponID primitive.ObjectID) (CouponSchema, error)
	FindByCodeCoupon(couponCode string) error
	InsertCouponConstraint(coupon CouponContraint) (CouponContraint, error)
	UpdateCouponConstraint(coupon CouponContraint) (CouponContraint, error)
	FindByIDCouponConstraint(constraintID primitive.ObjectID) (CouponContraint, error)
}

type RepositoryCol struct {
	couponCollection     *mongo.Collection
	constraintCollection *mongo.Collection
}

func newRepository() (RepositoryCol, error) {
	s, _ := db.GetMongo()
	couponCol := s.Database("discount").Collection("coupons")
	constraintCol := s.Database("discount").Collection("constraints")

	return RepositoryCol{
		couponCollection:     couponCol,
		constraintCollection: constraintCol,
	}, nil
}

func (d RepositoryCol) InsertCoupon(coupon CouponSchema) (CouponSchema, error) {

	_, err := d.couponCollection.InsertOne(context.Background(), coupon)
	if err != nil {
		return coupon, err
	}

	return coupon, nil
}

func (d RepositoryCol) UpdateCoupon(coupon CouponSchema) (CouponSchema, error) {

	return coupon, nil
}

func (d RepositoryCol) FindAllCoupons() ([]CouponSchema, error) {

	return nil, nil
}

func (d RepositoryCol) FindByIDCoupon(couponID primitive.ObjectID) (CouponSchema, error) {
	var coupon CouponSchema
	filter := bson.M{"_id": couponID}
	err := d.couponCollection.FindOne(context.Background(), filter).Decode(&coupon)
	return coupon, err
}

func (d RepositoryCol) FindByCodeCoupon(couponCode string) (CouponSchema, error) {
	var coupon CouponSchema
	filter := bson.D{{"code", couponCode}}
	res := d.couponCollection.FindOne(context.Background(), filter)
	err := res.Decode(&coupon)
	return coupon, err
}

func (d RepositoryCol) InsertCouponConstraint(constraint CouponContraint) (CouponContraint, error) {

	if _, err := d.constraintCollection.InsertOne(context.Background(), constraint); err != nil {
		return constraint, err
	}

	return constraint, nil
}

func (d RepositoryCol) UpdateCouponConstraint(coupon CouponContraint) (CouponContraint, error) {

	return coupon, nil
}

func (d RepositoryCol) FindByIDCouponConstraint(constraintID primitive.ObjectID) (CouponContraint, error) {
	var constraint CouponContraint
	filter := bson.M{"_id": constraintID}
	err := d.constraintCollection.FindOne(context.Background(), filter).Decode(&constraint)
	return constraint, err
}
