package main

import (
	"fmt"
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

	/* Step 1, make TCP connection to MQ */

	mqConn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		fatalError(err, "Failed to establish connection to RabbitMQ")
	}

	// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

	/* Step 2, make virtual connection to MQ by TCP */

	mqChannel, err = mqConn.Channel()
	if err != nil {
		fatalError(err, "Failed to establish connection to RabbitMQ")
	}

	// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

	/* Step 3, make exchange, bindings & queue*/
	/* not necessary to do it here, it can be done anywhere before */
	{
		/* 3.1 make exchange */

		err = mqChannel.ExchangeDeclare(
			"trial_exchange_1",  // name
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

		// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

		/* 3.2 make queue */

		mqQueue, err := mqChannel.QueueDeclare(
			"trial_queue_1", // name
			false,           // durable
			true,            // autoDelete
			true,            // exclusive
			false,           // noWait
			nil,             // args Table
		)
		if err != nil {
			fatalError(err, "Declare queue in RabbitMQ error")
		}

		fmt.Printf("mqQueue.Name: %s\n", mqQueue.Name)
		fmt.Printf("mqQueue.Consumers: %v\n", mqQueue.Consumers)
		fmt.Printf("mqQueue.Messages: %v\n", mqQueue.Messages)

		// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

		/* 3.3 make "Bindings" between "Exchange" and "Queue" under "Channel" */

		err = mqChannel.QueueBind(
			"trial_queue_1",      // queue name to bind to
			"trial_routingKey_1", // key, suppose be overruled by fanout type by exchange
			"trial_exchange_1",   // exchange instance to use
			false,                // noWait
			nil,                  // args Table
		)
		if err != nil {
			fatalError(err, "Bind queue to Exchange in RabbitMQ error")
		}
	}

	// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

	/* Step 4, set Consumer ["delivery" is type golang channel] */
	/* Create multiple consumer for One Queue */

	go CreateConsumer("trial_queue_1", "ConA")
	go CreateConsumer("trial_queue_1", "ConB")
	go CreateConsumer("trial_queue_1", "ConC")

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
