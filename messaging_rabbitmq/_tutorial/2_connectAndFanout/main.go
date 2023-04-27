package main

import (
	"log"

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

	{
		/* Make One Single Exchange */

		err = mqChannel.ExchangeDeclare(
			"trial_exchange_2",  // name
			amqp.ExchangeFanout, // type, way to distribute message
			false,               // durable
			true,                // auto-deleted
			false,               // internal
			false,               // no-wait
			nil,                 // arguments
		)
		if err != nil {
			fatalError(err, "Unable to create Exchange in RabbitMQ error")
		}

		// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

		/* Make Queue A */

		_, err := mqChannel.QueueDeclare(
			"trial_queue_2_A", // name
			false,             // durable
			true,              // autoDelete
			true,              // exclusive
			false,             // noWait
			nil,               // args Table
		)
		if err != nil {
			fatalError(err, "Declare queue in RabbitMQ error")
		}

		/* Make Queue B */

		_, err = mqChannel.QueueDeclare(
			"trial_queue_2_B", // name
			false,             // durable
			true,              // autoDelete
			true,              // exclusive
			false,             // noWait
			nil,               // args Table
		)
		if err != nil {
			fatalError(err, "Declare queue in RabbitMQ error")
		}

		// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

		/* Bind Queue A & B to Exchange */

		err = mqChannel.QueueBind(
			"trial_queue_2_A",    // queue name to bind to
			"trial_routingKey_2", // key, suppose be overruled by fanout type by exchange
			"trial_exchange_2",   // exchange instance to use
			false,                // noWait
			nil,                  // args Table
		)
		if err != nil {
			fatalError(err, "Declare queue in RabbitMQ error")
		}

		err = mqChannel.QueueBind(
			"trial_queue_2_B",    // queue name to bind to
			"trial_routingKey_2", // key, suppose be overruled by fanout type by exchange
			"trial_exchange_2",   // exchange instance to use
			false,                // noWait
			nil,                  // args Table
		)
		if err != nil {
			fatalError(err, "Declare queue in RabbitMQ error")
		}
	}

	// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

	/* Create consumer for each Queue */

	go CreateConsumer("trial_queue_2_A", "A")
	go CreateConsumer("trial_queue_2_B", "B")

	forever := make(chan interface{})
	log.Printf("[*] Waiting for messagees,To exit press CTRL+C")
	<-forever
}

func CreateConsumer(queueName, iden string) {

	delivery, err := mqChannel.Consume(
		queueName, // queue name to consume to
		"",        // consumer tag
		true,      // autoAck
		false,     // exclusive
		false,     // noLocal
		false,     // noWait
		nil,       // args Table
	)
	if err != nil {
		fatalError(err, "Create consumer "+iden+" to RabbitMQ error")
	}
	for d := range delivery {
		log.Printf("Received a message %s: %s", iden, d.Body)
	}
}
