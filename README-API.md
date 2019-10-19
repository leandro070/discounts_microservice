<a name="top"></a>
# Order Service en Java v0.1.0

Microservicio de Ordenes

- [RabbitMQ_GET](#rabbitmq_get)
	- [Logout de Usuarios](#logout-de-usuarios)
	


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


