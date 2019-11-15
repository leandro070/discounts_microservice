<a name="top"></a>
# Order Service en Java v0.1.0

Microservicio de Ordenes

- [Cupones](#cupones)
	- [Crear nuevo cupón](#crear-nuevo-cupón)
	- [Dar de baja un cupón](#dar-de-baja-un-cupón)
	- [Devuelve el cupón](#devuelve-el-cupón)
	
- [RabbitMQ](#rabbitmq)
	- [Get cupon](#get-cupon)
	- [Usar cupon](#usar-cupon)
	- [Validar cupon](#validar-cupon)
	
- [RabbitMQ_GET](#rabbitmq_get)
	- [Logout de Usuarios](#logout-de-usuarios)
	


# <a name='cupones'></a> Cupones

## <a name='crear-nuevo-cupón'></a> Crear nuevo cupón
[Back to top](#top)

<p>se encargará de crear un cupón y sus restricciones asociadas</p>

	POST /v1/coupons





### Request body Parameters

| Name     | Type       | Description                           |
|:---------|:-----------|:--------------------------------------|
|  description | String | <p>Nombre del descuento</p>|
|  amount | Number | **optional**<p>Monto fijo de descuento</p>|
|  percentaje | Number | **optional**<p>Porcentaje de descuento</p>|
|  coupon_type | String | <p>Tipo de descuento (&quot;fixed_amount&quot; o &quot;percentage&quot;)</p>|
|  constraint | Object | <p>Restricciones del cupón</p>|
| &nbsp;&nbsp;&nbsp;&nbsp; constraint.validity_from | Datetime | <p>Fecha de vigencia desde</p>|
| &nbsp;&nbsp;&nbsp;&nbsp; constraint.validity_to | Datetime | <p>Fecha de vigencia hasta</p>|
| &nbsp;&nbsp;&nbsp;&nbsp; constraint.total_usage | Number | <p>Total de usos que tiene el cupon actualmente</p>|
| &nbsp;&nbsp;&nbsp;&nbsp; constraint.max_usage | Number | <p>Máxima cantidad de usos posibles. 0 = Infinito</p>|
| &nbsp;&nbsp;&nbsp;&nbsp; constraint.max_amount | Number | <p>Si el descuento va por porcentaje entonces monto máximo establece un limite de descuento. 0 = Sin limite</p>|
| &nbsp;&nbsp;&nbsp;&nbsp; constraint.min_items | Number | <p>Permite que el descuento se aplique cuando hay mas de N items iguales. Ej: Un 2x1. 0 = Sin limites</p>|
| &nbsp;&nbsp;&nbsp;&nbsp; constraint.max_items | Number | <p>Permite que al descuento se aplique a una cierta cantidad de items. 0 = Sin limites</p>|
| &nbsp;&nbsp;&nbsp;&nbsp; constraint.combinable | Boolean | <p>Habilita que puedan combinarse con otros cupones.</p>|
### Examples

Body:

```
{
  "description": "Discount 1",
  "amount": 0,
  "percentage": null,
  "coupon_type": "fixed_amount",
  "constraint": {
      "validity_from": "2019-09-06T22:00:00.00-03:00",
      "validity_to": "2019-09-29T23:00:00.00-03:00",
      "total_usage": 1,
      "max_usage":20,
      "max_amount": 182372,
      "min_items": 2,
      "max_items": 5,
      "combinable": true
  }
}
```
Header Autorización

```
Authorization=bearer {token}
```




### Error Response

401 Unauthorized

```
HTTP/1.1 401 Unauthorized
```
400 Bad Request

```
HTTP/1.1 400 Bad Request
{
   "messages" : [
     {
       "path" : "{Nombre de la propiedad}",
       "message" : "{Motivo del error}"
     },
     ...
  ]
}
```
500 Server Error

```
HTTP/1.1 500 Internal Server Error
{
   "error" : "Not Found"
}
```
## <a name='dar-de-baja-un-cupón'></a> Dar de baja un cupón
[Back to top](#top)

<p>buscará un coupon y su restricción asociada y los dará de baja</p>

	DELETE /v1/coupons/:id



### Examples

Header Autorización

```
Authorization=bearer {token}
```


### Success Response

Respuesta

```
    HTTP/1.1 200 OK
	{
	  "id": "asd18cas18d15df1v6d"
	  "description": "Coupon $200",
	  "code": "ASD123"
	  "amount": 200,
	  "percentage": 0,
	  "coupon_type": "fixed_amount",
	  "constraint": {
	      "id": "asd18cas18d15df1v6d",
		  "validity_from": "2019-09-06T22:00:00.00-03:00",
	      "validity_to": "2019-09-29T23:00:00.00-03:00",
	      "total_usage": 1,
	      "max_usage": 20,
	      "max_amount": 182372,
	      "min_items": 2,
	      "max_items": 5,
	      "combinable": true
	  }
	}
```


### Error Response

401 Unauthorized

```
HTTP/1.1 401 Unauthorized
```
400 Bad Request

```
HTTP/1.1 400 Bad Request
{
   "messages" : [
     {
       "path" : "{Nombre de la propiedad}",
       "message" : "{Motivo del error}"
     },
     ...
  ]
}
```
500 Server Error

```
HTTP/1.1 500 Internal Server Error
{
   "error" : "Not Found"
}
```
## <a name='devuelve-el-cupón'></a> Devuelve el cupón
[Back to top](#top)

<p>buscará un coupon y su restricción asociada</p>

	POST /v1/coupons/:id



### Examples

Header Autorización

```
Authorization=bearer {token}
```


### Success Response

Respuesta

```
    HTTP/1.1 200 OK
	{
	  "id": "asd18cas18d15df1v6d"
	  "description": "Coupon $200",
	  "code": "ASD123"
	  "amount": 200,
	  "percentage": 0,
	  "coupon_type": "fixed_amount",
	  "constraint": {
	      "id": "asd18cas18d15df1v6d",
		  "validity_from": "2019-09-06T22:00:00.00-03:00",
	      "validity_to": "2019-09-29T23:00:00.00-03:00",
	      "total_usage": 1,
	      "max_usage": 20,
	      "max_amount": 182372,
	      "min_items": 2,
	      "max_items": 5,
	      "combinable": true
	  }
	}
```


### Error Response

401 Unauthorized

```
HTTP/1.1 401 Unauthorized
```
400 Bad Request

```
HTTP/1.1 400 Bad Request
{
   "messages" : [
     {
       "path" : "{Nombre de la propiedad}",
       "message" : "{Motivo del error}"
     },
     ...
  ]
}
```
500 Server Error

```
HTTP/1.1 500 Internal Server Error
{
   "error" : "Not Found"
}
```
# <a name='rabbitmq'></a> RabbitMQ

## <a name='get-cupon'></a> Get cupon
[Back to top](#top)

<p>DiscountService recibirá un codigo de cupon, validará que ese cupon esté vigente.</p>

	DIRECT discount_exchange/discount_requested/get_coupon_request





### Request body Parameters

| Name     | Type       | Description                           |
|:---------|:-----------|:--------------------------------------|
|  code | String | <p>Codigo del cupon</p>|
### Examples

Example usage

```
{
   "type": "get_coupon_request",
   "message": "{\"id\":\"5db36a4e80bea523b632258c\",\"description\":\"Discount 1\",\"amount\":0,\"percentage\":0,\"is_enable\":false,\"code\":\"UOATI8\",\"constraint\":{\"id\":\"5db36a4e80bea523b632258b\",\"validity_from\":\"2019-09-07T01:00:00Z\",\"validity_to\":\"2019-11-30T02:00:00Z\",\"total_usage\":19,\"max_usage\":20,\"max_amount\":182372,\"min_items\":2,\"max_items\":5,\"combinable\":true},\"coupon_type\":\"fixed_amount\"}"
}
```

### Param Examples

(json)
{RabbitMQ Message} Request-Example:

```
    RoutingKey: "discount_requested"
    Exchange: "discount_exchange"
    Queue: "discount_provider"
	   Type: "direct"
```
(json)
{RabbitMQ Message} Response-Example:

```
    RoutingKey: "discount_provided"
    Exchange: "discount_exchange"
    Queue: "discount_consumer"
	   Type: "direct"
```

### Success Response

Response-Mensaje

```
    {
       "type": "get_coupon_response",
       "message": "Coupon valid",
		  "success": "true"
    }
```


### Error Response

Error-Response-Mensaje

```
    {
       "type": "get_coupon_response",
       "message": "Max usage for coupon",
		  "success": "false"
    }
```
## <a name='usar-cupon'></a> Usar cupon
[Back to top](#top)

<p>DiscountService recibirá un codigo de cupon, validará que ese cupon esté vigente y luego aumentará la variable total_usage en 1. En caso de querer solo validar usar validate_request</p>

	DIRECT discount_exchange/discount_requested/use_coupon_request





### Request body Parameters

| Name     | Type       | Description                           |
|:---------|:-----------|:--------------------------------------|
|  code | String | <p>Codigo del cupon</p>|
|  items_to_apply | Number | <p>Cantidad de articulos que a los que se quiere aplicar</p>|
### Examples

Example usage

```
{
   "type": "use_coupon_request",
   "message": "{\"code\": \"UOATI8\", \"items_to_apply\": 2}"
}
```

### Param Examples

(json)
{RabbitMQ Message} Request-Example:

```
    RoutingKey: "discount_requested"
    Exchange: "discount_exchange"
    Queue: "discount_provider"
	   Type: "direct"
```
(json)
{RabbitMQ Message} Response-Example:

```
    RoutingKey: "discount_provided"
    Exchange: "discount_exchange"
    Queue: "discount_consumer"
	   Type: "direct"
```

### Success Response

Response-Mensaje

```
    {
       "type": "use_coupon_response",
       "message": "Success use cupon",
		  "success": "true"
    }
```


### Error Response

Error-Response-Mensaje

```
    {
       "type": "use_coupon_response",
       "message": "Coupon code empty",
		  "success": "false"
    }
```
## <a name='validar-cupon'></a> Validar cupon
[Back to top](#top)

<p>DiscountService recibirá un codigo de cupon, validará que ese cupon esté vigente.</p>

	DIRECT discount_exchange/discount_requested/validate_request





### Request body Parameters

| Name     | Type       | Description                           |
|:---------|:-----------|:--------------------------------------|
|  code | String | <p>Codigo del cupon</p>|
|  items_to_apply | Number | <p>Cantidad de articulos que a los que se quiere aplicar</p>|
### Examples

Example usage

```
{
   "type": "validate_request",
   "message": "{\"code\": \"UOATI8\", \"items_to_apply\": 2}"
}
```

### Param Examples

(json)
{RabbitMQ Message} Request-Example:

```
    RoutingKey: "discount_requested"
    Exchange: "discount_exchange"
    Queue: "discount_provider"
	   Type: "direct"
```
(json)
{RabbitMQ Message} Response-Example:

```
    RoutingKey: "discount_provided"
    Exchange: "discount_exchange"
    Queue: "discount_consumer"
	   Type: "direct"
```

### Success Response

Response-Mensaje

```
    {
       "type": "validate_response",
       "message": "Coupon valid",
		  "success": "true"
    }
```


### Error Response

Error-Response-Mensaje

```
    {
       "type": "validate_response",
       "message": "Max usage for coupon",
		  "success": "false"
    }
```
# <a name='rabbitmq_get'></a> RabbitMQ_GET

## <a name='logout-de-usuarios'></a> Logout de Usuarios
[Back to top](#top)

<p>Escucha de mensajes logout desde auth.</p>

	FANOUT auth/logout





### Success Response

Mensaje

```
{
   "type": "logout",
   "message": "{tokenId}"
}
```


