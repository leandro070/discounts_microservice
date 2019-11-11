package coupon

import (
	"encoding/json"
	"fmt"
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

/**
    * @api {DELETE} /v1/coupons/:id Dar de baja un cupón
    * @apiVersion 1.0.0
    * @apiName Delete coupon
	* @apiGroup Cupones
	* @apiDescription buscará un coupon y su restricción asociada y los dará de baja
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

/**
 * @api {direct} discount_exchange/discount_requested/use_coupon_request Usar cupon
 * @apiGroup RabbitMQ
 *
 * @apiDescription DiscountService recibirá un codigo de cupon, validará que ese cupon esté vigente y luego aumentará la variable total_usage en 1. En caso de querer solo validar usar validate_request
 *
 * @apiParamExample {RabbitMQ Message} Request-Example:
 *     RoutingKey: "discount_requested"
 *     Exchange: "discount_exchange"
 *     Queue: "discount_provider"
 *	   Type: "direct"
 *
 * @apiParam (Request body) {String} code Codigo del cupon
 * @apiParam (Request body) {Number} items_to_apply Cantidad de articulos que a los que se quiere aplicar
 *
 * @apiParamExample {RabbitMQ Message} Response-Example:
 *     RoutingKey: "discount_provided"
 *     Exchange: "discount_exchange"
 *     Queue: "discount_consumer"
 *	   Type: "direct"
 *
 * @apiExample {json} Example usage
 *     {
 *        "type": "use_coupon_request",
 *        "message": "{\"code\": \"UOATI8\", \"items_to_apply\": 2}"
 *     }
 *
 * @apiSuccessExample {json} Response-Mensaje
 *     {
 *        "type": "use_coupon_response",
 *        "message": "Success use cupon",
 *		  "success": "true"
 *     }
 *
 * @apiErrorExample {json} Error-Response-Mensaje
 *     {
 *        "type": "use_coupon_response",
 *        "message": "Coupon code empty",
 *		  "success": "false"
 *     }
 */
// UseCoupon se encargará de recibir un código de descuento y restricciones. Validará que el cupon sea aplicable y procederá a aplicarlo en caso de que sea real.
func UseCoupon(jsonData []byte) ([]byte, error) {
	var constraint struct {
		Code         string `json:"code"`
		ItemsToApply int    `json:"items_to_apply"`
	}

	err := json.Unmarshal(jsonData, &constraint)
	if err != nil {
		log.Println("ERROR - UseCoupon - unmarshall ", err.Error())
		return nil, err
	}

	serv, err := NewService()
	if err != nil {
		log.Println("ERROR - UseCoupon - create service", err.Error())
		return nil, err
	}

	if len(constraint.Code) == 0 {
		err := fmt.Errorf("Coupon code empty")
		log.Printf("ERROR - UseCoupon - %s", err.Error())
		return nil, err
	}

	log.Printf("INFO - UseCoupon - Coupon to use: %s", constraint.Code)
	err = serv.UseCoupon(constraint.Code, constraint.ItemsToApply)
	if err != nil {
		log.Println("ERROR - UseCoupon - Response:", err.Error())
		return nil, err
	}

	response := []byte("Success use cupon")

	return response, nil
}

/**
 * @api {direct} discount_exchange/discount_requested/validate_request Validar cupon
 * @apiGroup RabbitMQ
 *
 * @apiDescription DiscountService recibirá un codigo de cupon, validará que ese cupon esté vigente.
 *
 * @apiParamExample {RabbitMQ Message} Request-Example:
 *     RoutingKey: "discount_requested"
 *     Exchange: "discount_exchange"
 *     Queue: "discount_provider"
 *	   Type: "direct"
 *
 * @apiParam (Request body) {String} code Codigo del cupon
 * @apiParam (Request body) {Number} items_to_apply Cantidad de articulos que a los que se quiere aplicar
 *
 * @apiParamExample {RabbitMQ Message} Response-Example:
 *     RoutingKey: "discount_provided"
 *     Exchange: "discount_exchange"
 *     Queue: "discount_consumer"
 *	   Type: "direct"
 *
 * @apiExample {json} Example usage
 *     {
 *        "type": "validate_request",
 *        "message": "{\"code\": \"UOATI8\", \"items_to_apply\": 2}"
 *     }
 *
 * @apiSuccessExample {json} Response-Mensaje
 *     {
 *        "type": "validate_response",
 *        "message": "Coupon valid",
 *		  "success": "true"
 *     }
 *
 * @apiErrorExample {json} Error-Response-Mensaje
 *     {
 *        "type": "validate_response",
 *        "message": "Max usage for coupon",
 *		  "success": "false"
 *     }
 */
// ValidateCoupon se encargará de recibir un código de descuento y restricciones. Validará que el cupon sea aplicable.
func ValidateCoupon(jsonData []byte) ([]byte, error) {
	var constraint struct {
		Code         string `json:"code"`
		ItemsToApply int    `json:"items_to_apply"`
	}

	err := json.Unmarshal(jsonData, &constraint)
	if err != nil {
		log.Println("ERROR - ValidateCoupon - unmarshall ", err.Error())
		return nil, err
	}

	serv, err := NewService()
	if err != nil {
		log.Println("ERROR - ValidateCoupon - create service", err.Error())
		return nil, err
	}

	if len(constraint.Code) == 0 {
		err := fmt.Errorf("coupon code empty")
		log.Printf("ERROR - ValidateCoupon - %s", err.Error())
		return nil, err
	}

	log.Printf("INFO - ValidateCoupon - Coupon to use: %s", constraint.Code)
	err = serv.ValidateCoupon(constraint.Code, constraint.ItemsToApply)
	if err != nil {
		log.Println("ERROR - ValidateCoupon - Response:", err.Error())
		return nil, err
	}

	response := []byte("Coupon valid")

	return response, nil
}

/**
 * @api {direct} discount_exchange/discount_requested/get_coupon_request Get cupon
 * @apiGroup RabbitMQ
 *
 * @apiDescription DiscountService recibirá un codigo de cupon, validará que ese cupon esté vigente.
 *
 * @apiParamExample {RabbitMQ Message} Request-Example:
 *     RoutingKey: "discount_requested"
 *     Exchange: "discount_exchange"
 *     Queue: "discount_provider"
 *	   Type: "direct"
 *
 * @apiParam (Request body) {String} code Codigo del cupon
 *
 * @apiParamExample {RabbitMQ Message} Response-Example:
 *     RoutingKey: "discount_provided"
 *     Exchange: "discount_exchange"
 *     Queue: "discount_consumer"
 *	   Type: "direct"
 *
 * @apiExample {json} Example usage
 *     {
 *        "type": "get_coupon_request",
 *        "message": "{\"id\":\"5db36a4e80bea523b632258c\",\"description\":\"Discount 1\",\"amount\":0,\"percentage\":0,\"is_enable\":false,\"code\":\"UOATI8\",\"constraint\":{\"id\":\"5db36a4e80bea523b632258b\",\"validity_from\":\"2019-09-07T01:00:00Z\",\"validity_to\":\"2019-11-30T02:00:00Z\",\"total_usage\":19,\"max_usage\":20,\"max_amount\":182372,\"min_items\":2,\"max_items\":5,\"combinable\":true},\"coupon_type\":\"fixed_amount\"}"
 *     }
 *
 * @apiSuccessExample {json} Response-Mensaje
 *     {
 *        "type": "get_coupon_response",
 *        "message": "Coupon valid",
 *		  "success": "true"
 *     }
 *
 * @apiErrorExample {json} Error-Response-Mensaje
 *     {
 *        "type": "get_coupon_response",
 *        "message": "Max usage for coupon",
 *		  "success": "false"
 *     }
 */
// GetCouponByCode devolverá un cupon con su restriccion asociada
func GetCouponByCode(jsonData []byte) ([]byte, error) {
	var coupon struct {
		Code string `json:"code"`
	}

	err := json.Unmarshal(jsonData, &coupon)
	if err != nil {
		log.Println("ERROR - GetCouponByCode - unmarshall ", err.Error())
		return nil, err
	}

	if len(coupon.Code) == 0 {
		err := fmt.Errorf("coupon code empty")
		log.Printf("ERROR - GetCouponByCode - %s", err.Error())
		return nil, err
	}

	serv, err := NewService()
	if err != nil {
		log.Println("ERROR - GetCouponByCode - create service", err.Error())
		return nil, err
	}

	log.Printf("INFO - GetCouponByCode - Coupon to use: %s", coupon.Code)
	result, err := serv.GetCouponByCode(coupon.Code)
	if err != nil {
		log.Println("ERROR - GetCouponByCode - Response:", err.Error())
		return nil, err
	}

	response, err := json.Marshal(result)
	if err != nil {
		log.Println("ERROR - GetCouponByCode - Marshall:", err.Error())
		return nil, err
	}

	return response, nil
}
