package coupon

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/leandro070/discounts_microservice/utils/errors"
)

/**
    * @api {post} /v1/coupons Crear nuevo cupón
    * @apiVersion 1.0.0
    * @apiName Create coupon
	* @apiGroup Cupones
	* @apiDescription se encargará de crear un cupón y sus restricciones asociadas
    *
    * @apiParam (Request body) {String} description Nombre del descuento
    * @apiParam (Request body) {Number} [amount] Monto fijo de descuento
    * @apiParam (Request body) {Number} [percentaje] Porcentaje de descuento
    * @apiParam (Request body) {String} coupon_type Tipo de descuento ("fixed_amount" o "percentage")
    * @apiParam (Request body) {Object} constraint Restricciones del cupón
    * @apiParam (Request body) {Datetime} constraint.validity_from Fecha de vigencia desde
    * @apiParam (Request body) {Datetime} constraint.validity_to Fecha de vigencia hasta
	* @apiParam (Request body) {Number} constraint.total_usage Total de usos que tiene el cupon actualmente
	* @apiParam (Request body) {Number} constraint.max_usage Máxima cantidad de usos posibles. 0 = Infinito
	* @apiParam (Request body) {Number} constraint.max_amount Si el descuento va por porcentaje entonces monto máximo establece un limite de descuento. 0 = Sin limite
	* @apiParam (Request body) {Number} constraint.min_items Permite que el descuento se aplique cuando hay mas de N items iguales. Ej: Un 2x1. 0 = Sin limites
	* @apiParam (Request body) {Number} constraint.max_items Permite que al descuento se aplique a una cierta cantidad de items. 0 = Sin limites
	* @apiParam (Request body) {Boolean} constraint.combinable Habilita que puedan combinarse con otros cupones.
	*
    * @apiExample {json} Body:
	* 	{
	* 	  "description": "Discount 1",
	* 	  "amount": 0,
	* 	  "percentage": null,
	* 	  "coupon_type": "fixed_amount",
	* 	  "constraint": {
	* 	      "validity_from": "2019-09-06T22:00:00.00-03:00",
	* 	      "validity_to": "2019-09-29T23:00:00.00-03:00",
	* 	      "total_usage": 1,
	* 	      "max_usage":20,
	* 	      "max_amount": 182372,
	* 	      "min_items": 2,
	* 	      "max_items": 5,
	* 	      "combinable": true
	* 	  }
	* 	}
    *
 * @apiUse AuthHeader
 * @apiUse ParamValidationErrors
 * @apiUse OtherErrors
*/
// NewCoupon se encargará de crear un cupón y sus restricciones asociadas
func NewCoupon(c *gin.Context) {

	var coupon NewCouponRequest
	err := c.ShouldBindJSON(&coupon)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	couponService, err := NewService()
	if err != nil {
		errors.Handle(c, err)
		return
	}

	res, err := couponService.NewCoupon(&coupon)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.JSON(200, gin.H{
		"result": res,
	})
}

/**
    * @api {post} /v1/coupons/:id Devuelve el cupón
    * @apiVersion 1.0.0
    * @apiName Get coupon
	* @apiGroup Cupones
	* @apiDescription buscará un coupon y su restricción asociada
    *
    * @apiParam (Request body) {String} description Nombre del descuento
 	*
    * @apiSuccessExample {json} Respuesta
 	*     HTTP/1.1 200 OK
	* 	{
	* 	  "id": "asd18cas18d15df1v6d"
	* 	  "description": "Coupon $200",
	*	  "code": "ASD123"
	* 	  "amount": 200,
	* 	  "percentage": 0,
	* 	  "coupon_type": "fixed_amount",
	* 	  "constraint": {
	* 	      "id": "asd18cas18d15df1v6d",
	*		  "validity_from": "2019-09-06T22:00:00.00-03:00",
	* 	      "validity_to": "2019-09-29T23:00:00.00-03:00",
	* 	      "total_usage": 1,
	* 	      "max_usage": 20,
	* 	      "max_amount": 182372,
	* 	      "min_items": 2,
	* 	      "max_items": 5,
	* 	      "combinable": true
	* 	  }
	* 	}
    *
 	* @apiUse AuthHeader
 	* @apiUse ParamValidationErrors
 	* @apiUse OtherErrors
*/
// GetCoupon se encargará de recibir un código de descuento, validar la existencia y vigencia del cupón.
func GetCoupon(c *gin.Context) {

	couponID := c.Param("id")

	couponService, err := NewService()
	if err != nil {
		errors.Handle(c, err)
		return
	}

	res, err := couponService.GetCoupon(couponID)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.JSON(200, gin.H{
		"result": res,
	})
}

// AnnulCoupon se encargará de recibir un código de descuento y darlo de baja.
func AnnulCoupon(c *gin.Context) {

	couponID := c.Param("id")

	couponService, err := NewService()
	if err != nil {
		errors.Handle(c, err)
		return
	}

	err = couponService.AnnulCoupon(couponID)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.JSON(200, gin.H{
		"success": true,
	})
}

func UseCoupon(code string) ([]byte, error) {

	serv, err := NewService()
	if err != nil {
		return nil, err
	}

	err = serv.UseCoupon(code)
	if err != nil {
		return nil, err
	}

	result := gin.H{"type": "use_coupon_response", "message": "true"}

	response, err := json.Marshal(result)
	if err != nil {
		log.Printf("ERROR: fail marshal: %s", result)
	}

	return response, nil
}
