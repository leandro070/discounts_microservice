package coupon

import (
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	validator "gopkg.in/go-playground/validator.v9"
)

// ServiceCupon contendrá instancias a repositorio
type ServiceCupon struct {
	repo RepositoryCol
}

// Service es la interfaz que contenga todas las acciones a realizar
type Service interface {
	NewCoupon(coupon *NewCouponRequest) (CouponConstraintResponse, error)
	GetCoupon(couponID string) (CouponResponse, error)
}

// NewService retorna una nueva instancia del servicio
func NewService() (ServiceCupon, error) {
	var service ServiceCupon
	repo, err := newRepository()
	if err != nil {
		return service, err
	}
	service.repo = repo
	return service, nil
}

// NewCoupon es el servicio que creará un nuevo cupon
func (s ServiceCupon) NewCoupon(req *NewCouponRequest) (CouponResponse, error) {
	var res CouponResponse

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		return res, err
	}

	code := newCode(&s)

	constraint := CouponContraint{
		ID:           primitive.NewObjectID(),
		ValidityFrom: req.Constraint.ValidityFrom,
		ValidityTo:   req.Constraint.ValidityTo,
		TotalUsage:   req.Constraint.TotalUsage,
		MaxUsage:     req.Constraint.MaxUsage,
		MaxAmount:    req.Constraint.MaxAmount,
		MaxItems:     req.Constraint.MaxItems,
		MinItems:     req.Constraint.MinItems,
		Combinable:   req.Constraint.Combinable,
		IsEnable:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	constraint, err = s.repo.InsertCouponConstraint(constraint)
	if err != nil {
		return res, err
	}

	coupon := CouponSchema{
		ID:           primitive.NewObjectID(),
		Description:  req.Description,
		Code:         code,
		Amount:       req.Amount,
		Percentage:   req.Percentage,
		IsEnable:     true,
		ConstraintID: constraint.ID,
		CouponType:   req.CouponType,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	coupon, err = s.repo.InsertCoupon(coupon)
	if err != nil {
		return res, err
	}

	res = CouponResponse{
		ID:          coupon.ID,
		Description: coupon.Description,
		Code:        coupon.Code,
		Amount:      coupon.Amount,
		Percentage:  coupon.Percentage,
		IsEnable:    coupon.IsEnable,
		CouponType:  req.CouponType,
		Constraint: CouponConstraintResponse{
			ID:           constraint.ID,
			ValidityFrom: constraint.ValidityFrom,
			ValidityTo:   constraint.ValidityTo,
			TotalUsage:   constraint.TotalUsage,
			MaxUsage:     constraint.MaxUsage,
			MaxAmount:    constraint.MaxAmount,
			MaxItems:     constraint.MaxItems,
			MinItems:     constraint.MinItems,
			Combinable:   constraint.Combinable,
		},
	}

	return res, nil
}

func newCode(s *ServiceCupon) string {
	code := codeRandomGenerator(6)
	ok := validateCode(s, code)
	for ok == false {
		code = codeRandomGenerator(6)
		ok = validateCode(s, code)
	}
	return code
}

func codeRandomGenerator(n int) string {
	rand.Seed(time.Now().UnixNano())
	letterRunes := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// GetCoupon retorna un cupon segun el id enviado
func (s ServiceCupon) GetCoupon(couponID string) (CouponResponse, error) {
	var resp CouponResponse
	id, err := primitive.ObjectIDFromHex(couponID)
	if err != nil {
		return resp, err
	}
	coupon, err := s.repo.FindByIDCoupon(id)
	if err != nil {
		return resp, err
	}

	constraint, err := s.repo.FindByIDCouponConstraint(coupon.ConstraintID)
	if err != nil {
		return resp, err
	}

	resp = CouponResponse{
		ID:          coupon.ID,
		Description: coupon.Description,
		Amount:      coupon.Amount,
		IsEnable:    coupon.IsEnable,
		Code:        coupon.Code,
		Percentage:  coupon.Percentage,
		Constraint: CouponConstraintResponse{
			ID:           constraint.ID,
			ValidityFrom: constraint.ValidityFrom,
			ValidityTo:   constraint.ValidityTo,
			TotalUsage:   constraint.TotalUsage,
			MaxAmount:    constraint.MaxAmount,
			MaxItems:     constraint.MaxItems,
			MaxUsage:     constraint.MaxUsage,
			MinItems:     constraint.MinItems,
			Combinable:   constraint.Combinable,
		},
		CouponType: coupon.CouponType,
	}

	return resp, nil
}

/*
	Al momento de dar de baja un coupon debo comunicarme con ORDER para saber si el cupon está en uso
	en alguna ORDER. is_enable en false si o si...
*/
// AnnulCoupon da de baja un cupon
func (s ServiceCupon) AnnulCoupon(couponID string) error {

	id, err := primitive.ObjectIDFromHex(couponID)
	if err != nil {
		return err
	}

	coupon, err := s.repo.FindByIDCoupon(id)
	if err != nil {
		return err
	}

	err = s.repo.AnnulCoupon(coupon.ID)
	if err != nil {
		return err
	}

	err = couponDisable(couponID)
	if err != nil {
		return err
	}

	return nil
}

func (s ServiceCupon) UseCoupon(couponCode string) error {
	return nil
}
