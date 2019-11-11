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
	Success bool   `json:"success"`
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
		"discount_fanout_exchange", // name
		"fanout",                   // type
		false,                      // durable
		false,                      // auto-deleted
		false,                      // internal
		false,                      // no-wait
		nil,                        // arguments
	)
	if err != nil {
		return err
	}

	body, err := json.Marshal(send)
	if err != nil {
		return err
	}

	err = channel.Publish(
		"discount_fanout_exchange", // exchange
		"",                         // routing key
		false,                      // mandatory
		false,                      // immediate
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

	err = amqpChannel.ExchangeDeclare(
		"discount_exchange", // name
		"direct",            // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		log.Printf("ERROR: fail declare exchange : %s", err.Error())
		return err
	}

	_, err = amqpChannel.QueueDeclare(
		"discount_provider", //name
		true,                //durable
		false,               //autodeleted
		false,               //exclusive
		false,               //noWait
		nil,                 //args
	)
	if err != nil {
		log.Printf("ERROR: fail create queue discount_provider: %s", err.Error())
	}

	err = amqpChannel.QueueBind(
		"discount_provider",
		"discount_requested",
		"discount_exchange",
		false,
		nil,
	)
	if err != nil {
		log.Printf("ERROR: fail bind queue discount_provider: %s", err.Error())
	}

	_, err = amqpChannel.QueueDeclare(
		"discount_consumer", //name
		true,                //durable
		false,               //autodeleted
		false,               //exclusive
		false,               //noWait
		nil,                 //args
	)
	if err != nil {
		log.Printf("ERROR: fail create queue discount_consumer: %s", err.Error())
	}
	err = amqpChannel.QueueBind(
		"discount_consumer",
		"discount_provided",
		"discount_exchange",
		false,
		nil,
	)
	if err != nil {
		log.Printf("ERROR: fail bind queue discount_consumer: %s", err.Error())
	}

	// channel
	msgChannel, err := amqpChannel.Consume(
		"discount_provider", // queue
		"",                  // consumer
		true,                // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)
	if err != nil {
		log.Printf("ERROR: fail consume queue discount_provider: %s", err.Error())
		return err
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

				log.Printf("INFO: coupon type %s", msg.Type)
				msgResp, err := UseCoupon([]byte(msg.Message))
				var response []byte
				if err != nil {
					response, err = makeResponse("use_coupon_response", []byte(err.Error()), false)
				} else {
					response, err = makeResponse("use_coupon_response", msgResp, true)
				}
				if err != nil {
					channel = nil
					return
				}

				err = amqpChannel.Publish(
					"discount_exchange", // exchange
					"discount_provided", // routing key
					false,               // mandatory
					false,               // immediate
					amqp.Publishing{
						ContentType: "application/json",
						Body:        response,
					},
				)
				if err != nil {
					channel = nil
					return
				}
				log.Println("Mensaje reenviado a:", "[REPLY_TO]", rabbitMsg.ReplyTo, "[CORRELATION_ID]", rabbitMsg.CorrelationId)

			case "validate_request":

				log.Printf("INFO: coupon type %s", msg.Type)
				msgResp, err := ValidateCoupon([]byte(msg.Message))
				var response []byte
				if err != nil {
					response, err = makeResponse("validate_response", []byte(err.Error()), false)
				} else {
					response, err = makeResponse("validate_response", msgResp, true)
				}
				if err != nil {
					channel = nil
					return
				}

				err = amqpChannel.Publish(
					"discount_exchange", // exchange
					"discount_provided", // routing key
					false,               // mandatory
					false,               // immediate
					amqp.Publishing{
						ContentType: "application/json",
						Body:        response,
					},
				)
				if err != nil {
					channel = nil
					return
				}
				log.Println("Mensaje reenviado a:", "[REPLY_TO]", rabbitMsg.ReplyTo, "[CORRELATION_ID]", rabbitMsg.CorrelationId)

			case "get_coupon_request":

				log.Printf("INFO: coupon type %s", msg.Type)
				msgResp, err := GetCouponByCode([]byte(msg.Message))
				var response []byte
				if err != nil {
					response, err = makeResponse("get_coupon_response", []byte(err.Error()), false)
				} else {
					response, err = makeResponse("get_coupon_response", msgResp, true)
				}
				if err != nil {
					channel = nil
					return
				}

				err = amqpChannel.Publish(
					"discount_exchange", // exchange
					"discount_provided", // routing key
					false,               // mandatory
					false,               // immediate
					amqp.Publishing{
						ContentType: "application/json",
						Body:        response,
					},
				)
				if err != nil {
					channel = nil
					return
				}
				log.Println("Mensaje reenviado a:", "[REPLY_TO]", rabbitMsg.ReplyTo, "[CORRELATION_ID]", rabbitMsg.CorrelationId)

			}
		}
	}()

	log.Print("Closed connection: ", <-conn.NotifyClose(make(chan *amqp.Error)))

	return nil
}

func makeResponse(t string, data []byte, success bool) ([]byte, error) {
	msg := message{
		Message: string(data),
		Type:    t,
		Success: success,
	}
	resp, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
