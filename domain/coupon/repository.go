package coupon

import (
	"context"
	"log"
	"time"

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
	FindByCodeContraint(couponCode string) (CouponContraint, error)
	FindByIDCouponConstraint(constraintID primitive.ObjectID) (CouponContraint, error)
	AnnulCoupon(couponID primitive.ObjectID) error
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

/*#############################################################################
############################ COUPONS ######################################
#############################################################################*/

func (d RepositoryCol) InsertCoupon(coupon CouponSchema) (CouponSchema, error) {

	_, err := d.couponCollection.InsertOne(context.Background(), coupon)
	if err != nil {
		return coupon, err
	}

	return coupon, nil
}

func (d RepositoryCol) UpdateCoupon(coupon CouponSchema) (CouponSchema, error) {

	coupon.UpdatedAt = time.Now()

	filter := bson.M{"_id": bson.M{"$eq": coupon.ID}}

	_, err := d.couponCollection.UpdateOne(context.Background(), filter, coupon)

	if err != nil {
		return coupon, err
	}

	return coupon, nil
}

func (d RepositoryCol) AnnulCoupon(couponID primitive.ObjectID) error {
	var coupon CouponSchema
	filter := bson.M{"_id": bson.M{"$eq": couponID}}

	update := bson.M{
		"$set": bson.M{
			"is_enable":  false,
			"updated_at": time.Now(),
		},
	}

	res := d.couponCollection.FindOneAndUpdate(context.Background(), filter, update, nil)
	if res.Err() != nil {
		return res.Err()
	}

	err := res.Decode(&coupon)
	if err != nil {
		return err
	}

	err = d.AnnulContraint(coupon.ConstraintID)
	if err != nil {
		return err
	}

	return nil
}

func (d RepositoryCol) FindAllCoupons() ([]CouponSchema, error) {

	return nil, nil
}

func (d RepositoryCol) FindByIDCoupon(couponID primitive.ObjectID) (CouponSchema, error) {
	var coupon CouponSchema
	filter := bson.M{"_id": bson.M{"$eq": couponID}}
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

func (d RepositoryCol) FindByCodeContraint(couponCode string) (CouponContraint, error) {
	var constraint CouponContraint
	// pipeline := []bson.M{
	// 	bson.M{"$lookup": bson.M{
	// 		"from":         "cupons",
	// 		"localField":   "_id",
	// 		"foreignField": "constraint_id",
	// 		"as":           "coupon",
	// 		"pipeline": []bson.M{
	// 			"$match": bson.M{
	// 				"$expr": bson.M{
	// 					"$eq": bson.M{"code": couponCode, "is_enable": true},
	// 				},
	// 			}},
	// 	}},
	// }
	// cur, err := d.couponCollection.FindOne(context.Background(), pipeline)
	// if err != nil {
	// 	log.Fatal("[error]", err.Error())
	// 	return constraints, err
	// }
	// defer cur.Close(context.Background())
	// log.Println("sin error")
	// elem := &bson.D{}
	// cur.Decode(elem)
	// log.Println("[decoded]", elem)
	// for cur.Next(context.Background()) {
	// 	elem := &bson.D{}
	// 	if err := cur.Decode(elem); err != nil {
	// 		log.Fatal(err.Error())
	// 	}
	// 	log.Println(elem)
	// }

	coupon, err := d.FindByCodeCoupon(couponCode)
	if err != nil {
		return constraint, err
	}
	constraint, err = d.FindByIDCouponConstraint(coupon.ConstraintID)
	if err != nil {
		return constraint, err
	}
	log.Println("[CUPON]", coupon)
	log.Println("[REST]", constraint)
	return constraint, nil
}

/*#############################################################################
############################ CONSTRAINTS ######################################
#############################################################################*/

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
	filter := bson.M{"_id": bson.M{"$eq": constraintID}}
	err := d.constraintCollection.FindOne(context.Background(), filter).Decode(&constraint)
	return constraint, err
}

func (d RepositoryCol) AnnulContraint(constraintID primitive.ObjectID) error {

	filter := bson.M{"_id": bson.M{"$eq": constraintID}}

	update := bson.M{
		"$set": bson.M{
			"is_enable":  false,
			"updated_at": time.Now(),
		},
	}

	_, err := d.constraintCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
