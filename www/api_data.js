define({ "api": [
  {
    "type": "fanout",
    "url": "auth/logout",
    "title": "Logout de Usuarios",
    "group": "RabbitMQ_GET",
    "description": "<p>Escucha de mensajes logout desde auth.</p>",
    "success": {
      "examples": [
        {
          "title": "Mensaje",
          "content": "{\n   \"type\": \"logout\",\n   \"message\": \"{tokenId}\"\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./gateway/rabbit/rabbit.go",
    "groupTitle": "RabbitMQ_GET",
    "name": "FanoutAuthLogout"
  }
] });
