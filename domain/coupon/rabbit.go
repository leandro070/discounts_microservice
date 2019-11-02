package coupon

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/leandro070/discounts_microservice/utils/env"
	"github.com/leandro070/discounts_microservice/utils/errors"
	"github.com/leandro070/discounts_microservice/utils/security"
	"github.com/streadway/amqp"
)

// ErrChannelNotInitialized Rabbit channel could not be initialized
var ErrChannelNotInitialized = errors.NewCustom(400, "Channel not initialized")

var channel *amqp.Channel

type message struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// Init se queda escuchando broadcasts de logout
func RabbitInit() {
	go func() {
		for {
			listenLogout()
			fmt.Println("RabbitMQ conectando en 5 segundos.")
			time.Sleep(5 * time.Second)
		}
	}()
}

func getChannel() (*amqp.Channel, error) {
	if channel == nil {
		conn, err := amqp.Dial(env.Get().RabbitURL)
		if err != nil {
			return nil, err
		}

		ch, err := conn.Channel()
		if err != nil {
			return nil, err
		}
		channel = ch
	}
	if channel == nil {
		return nil, ErrChannelNotInitialized
	}
	return channel, nil
}

/**
 * @api {fanout} auth/logout Logout de Usuarios
 * @apiGroup RabbitMQ GET
 *
 * @apiDescription Escucha de mensajes logout desde auth.
 *
 * @apiSuccessExample {json} Mensaje
 *     {
 *        "type": "logout",
 *        "message": "{tokenId}"
 *     }
 */
func listenLogout() error {
	conn, err := amqp.Dial(env.Get().RabbitURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	chn, err := conn.Channel()
	if err != nil {
		return err
	}
	defer chn.Close()

	err = chn.ExchangeDeclare(
		"auth",   // name
		"fanout", // type
		false,    // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return err
	}

	queue, err := chn.QueueDeclare(
		"auth", // name
		false,  // durable
		false,  // delete when unused
		true,   // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		return err
	}

	err = chn.QueueBind(
		queue.Name, // queue name
		"",         // routing key
		"auth",     // exchange
		false,
		nil)
	if err != nil {
		return err
	}

	mgs, err := chn.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return err
	}

	fmt.Println("RabbitMQ conectado")

	go func() {
		for d := range mgs {
			log.Output(1, "Mensage recibido")
			newMessage := &message{}
			err = json.Unmarshal(d.Body, newMessage)
			if err == nil {
				if newMessage.Type == "logout" {
					security.Invalidate(newMessage.Message)
				}
			}
		}
	}()

	fmt.Print("Closed connection: ", <-conn.NotifyClose(make(chan *amqp.Error)))

	return nil
}

func CouponDisable(couponID string) error {
	send := message{
		Type:    "coupon_disabled",
		Message: couponID,
	}

	channel, err := getChannel()
	if err != nil {
		channel = nil
		return err
	}

	err = channel.ExchangeDeclare(
		"discount", // name
		"fanout",   // type
		false,      // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)

	body, err := json.Marshal(send)
	if err != nil {
		return err
	}

	err = channel.Publish(
		"auth", // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			Body: []byte(body),
		})
	if err != nil {
		channel = nil
		return err
	}

	log.Output(1, "Rabbit discunt disabled enviado")

	return nil
}
