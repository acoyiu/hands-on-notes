package main

import (
	"log"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

var mqConn *amqp.Connection
var mqChannel *amqp.Channel

func fatalError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s:%s", msg, err)
	}
}

func main() {
	var err error

	mqConn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		fatalError(err, "Failed to establish connection to RabbitMQ")
	}

	// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

	mqChannel, err = mqConn.Channel()
	if err != nil {
		fatalError(err, "Failed to establish connection to RabbitMQ")
	}

	// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

	r := gin.Default()
	r.GET("/send", func(c *gin.Context) {

		exchange := c.Query("exchange")
		if exchange == "" {
			c.JSON(200, gin.H{
				"message": "Missing exchange",
			})
			return
		}

		topic := c.Query("topic")
		if topic == "" {
			c.JSON(200, gin.H{
				"message": "Missing topic",
			})
			return
		}

		println("Publising to: " + exchange + topic)

		/*
			1: Publishing Error: Channel.NotifyReturn to handle any undeliverable message when calling publish

			2: mandatory: Return Error when no queue receive the message

			3:
		*/

		err = mqChannel.Publish(
			exchange, // exchange name   - instance to use
			topic,    // routingKey  	 - topic to use
			false,    // mandatory	     - Return Error when No queue receive the message
			false,    // immediate       - Return Error when No matched queue is ready to accept
			amqp.Publishing{
				Body: []byte("the body of msg published"),
				// Other Common params:
				// 	ContentType
				// 	Priority
				// 	CorrelationId
				// 	ReplyTo
				// 	Expiration
				// 	MessageId
				// 	Timestamp
				// 	Type
				// 	UserId
				// 	AppId
			},
		)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
		}

		c.JSON(200, gin.H{
			"message": true,
		})
	})
	r.Run(":4000")

}
