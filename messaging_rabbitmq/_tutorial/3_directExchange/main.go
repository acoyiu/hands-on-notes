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
			"trial_exchange_3",  // name
			amqp.ExchangeDirect, // type, way to distribute message
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

		mqChannel.QueueDeclare(
			"trial_queue_3_A", // name
			false,             // durable
			true,              // autoDelete
			true,              // exclusive
			false,             // noWait
			nil,               // args Table
		)

		/* Make Queue B */

		mqChannel.QueueDeclare(
			"trial_queue_3_B", // name
			false,             // durable
			true,              // autoDelete
			true,              // exclusive
			false,             // noWait
			nil,               // args Table
		)

		/* Make Queue C, For Error in purpose */
		/* Declare Queue C for using RoutingKey B */

		mqChannel.QueueDeclare(
			"trial_queue_3_C", // name
			false,             // durable
			true,              // autoDelete
			true,              // exclusive
			false,             // noWait
			nil,               // args Table
		)

		// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

		/* Bind Queues to Exchange */

		/* Each Queue is Depeneding on SAME NAME "Routing Key" */

		mqChannel.QueueBind(
			"trial_queue_3_A",      // queue name to bind to
			"trial_routingKey_3_A", // Routing key, suppose be overruled by fanout type by exchange
			"trial_exchange_3",     // exchange instance to use
			false,                  // noWait
			nil,                    // args Table
		)

		mqChannel.QueueBind(
			"trial_queue_3_B",      // queue name to bind to
			"trial_routingKey_3_B", // Routing key, suppose be overruled by fanout type by exchange
			"trial_exchange_3",     // exchange instance to use
			false,                  // noWait
			nil,                    // args Table
		)

		mqChannel.QueueBind(
			"trial_queue_3_C",      // queue name to bind to
			"trial_routingKey_3_B", // Routing key, suppose be overruled by fanout type by exchange
			"trial_exchange_3",     // exchange instance to use
			false,                  // noWait
			nil,                    // args Table
		)
	}

	// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

	/* Create consumer for each Queue */

	go CreateConsumer("trial_queue_3_A", "A")
	go CreateConsumer("trial_queue_3_B", "B")
	go CreateConsumer("trial_queue_3_C", "C")

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
