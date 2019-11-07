package coupon

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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

			log.Printf("RabbitMQ connecting in 5 seconds.")
			time.Sleep(5 * time.Second)
		}
	}()

	go func() {
		for {
			listenCouponsEvents()

			log.Printf("RabbitMQ connecting in 5 seconds")
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

	log.Printf("[*] Waiting for messages LOGOUT.")

	go func() {
		for d := range mgs {
			log.Printf("Mensage recibido")
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

func couponDisable(couponID string) error {
	send := message{
		Type:    "couponID",
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
		"discount", // exchange
		"",         // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Body: []byte(body),
		})
	if err != nil {
		channel = nil
		return err
	}

	log.Printf("Rabbit discount disabled enviado")

	return nil
}

func listenCouponsEvents() error {
	conn, err := amqp.Dial(env.Get().RabbitURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	if err != nil {
		log.Printf("ERROR: fail create channel: %s", err.Error())
		os.Exit(1)
	}
	defer amqpChannel.Close()

	// create queue
	queue, err := amqpChannel.QueueDeclare(
		"discount", // channelname
		false,      // durable
		false,      // delete when unused
		true,       // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Printf("ERROR: fail create queue: %s", err.Error())
		os.Exit(1)
	}

	err = amqpChannel.Qos(
		1,     // recuento de captación previa
		0,     // tamaño de captación previa
		false, // global
	)
	if err != nil {
		log.Printf("ERROR: Error al establecer QoS: %s", err.Error())
		os.Exit(1)
	}

	// channel
	msgChannel, err := amqpChannel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.Printf("ERROR: fail create channel: %s", err.Error())
		os.Exit(1)
	}

	log.Printf("[*] Waiting for messages USE COUPON.")

	go func() {
		for rabbitMsg := range msgChannel {
			var msg message
			err := json.Unmarshal(rabbitMsg.Body, &msg)
			if err != nil {
				log.Printf("ERROR: fail unmarshl: %s", rabbitMsg.Body)
				return
			}
			log.Printf("INFO: received msg: %v", msg)
			switch msg.Type {
			case "use_coupon_request":
				//TODO: Hacer lo que hay que hacer

				response, err := UseCoupon(msg.Message)
				if err != nil {
					return
				}

				err = amqpChannel.Publish(
					"",                // exchange
					rabbitMsg.ReplyTo, // routing key
					false,             // mandatory
					false,             // immediate
					amqp.Publishing{
						ContentType:   "application/json",
						Body:          response,
						CorrelationId: rabbitMsg.CorrelationId,
					})
				if err != nil {
					channel = nil
					return
				}
				log.Println("Mensaje reenviado a:", "[REPLY_TO]", rabbitMsg.ReplyTo, "[CORRELATION_ID]", rabbitMsg.CorrelationId)

				rabbitMsg.Ack(true)
			}
		}
	}()

	log.Print("Closed connection: ", <-conn.NotifyClose(make(chan *amqp.Error)))

	return nil
}

// {"type": "use_coupon", "message": "UOATI8"}

// reply := spec.CreateDocumentReply{
// 	Uid:    docMsg.Uid,
// 	Status: "Created",
// }
// msg := RabbitMsg{
// 	QueueName: docMsg.ReplyTo,
// 	Reply:     reply,
// }
