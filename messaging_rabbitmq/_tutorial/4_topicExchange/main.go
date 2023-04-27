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
			"trial_exchange_4",
			amqp.ExchangeTopic,
			false,
			true,
			false,
			false,
			nil,
		)
		if err != nil {
			fatalError(err, "Unable to create Exchange in RabbitMQ error")
		}

		// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

		/* Make Queue A */

		mqChannel.QueueDeclare(
			"trial_queue_4_A",
			false,
			true,
			true,
			false,
			nil,
		)

		/* Make Queue B */

		mqChannel.QueueDeclare(
			"trial_queue_4_B",
			false,
			true,
			true,
			false,
			nil,
		)

		/* Make Queue C */

		mqChannel.QueueDeclare(
			"trial_queue_4_C",
			false,
			true,
			true,
			false,
			nil,
		)

		/* Make Queue D */

		mqChannel.QueueDeclare(
			"trial_queue_4_D",
			false,
			true,
			true,
			false,
			nil,
		)

		// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

		/* Bind Queues to Exchange */

		/* QA QB will using routingKey ..01 */
		/* QC QD will using routingKey ..02 */

		mqChannel.QueueBind(
			"trial_queue_4_A",
			"usa.#",
			"trial_exchange_4",
			false,
			nil,
		)

		mqChannel.QueueBind(
			"trial_queue_4_B",
			"#.news",
			"trial_exchange_4",
			false,
			nil,
		)

		mqChannel.QueueBind(
			"trial_queue_4_C",
			"#.weather",
			"trial_exchange_4",
			false,
			nil,
		)

		mqChannel.QueueBind(
			"trial_queue_4_D",
			"europe.#",
			"trial_exchange_4",
			false,
			nil,
		)
	}

	// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

	/* Create consumer for each Queue */

	go CreateConsumer("trial_queue_4_A", "A")
	go CreateConsumer("trial_queue_4_B", "B")
	go CreateConsumer("trial_queue_4_C", "C")
	go CreateConsumer("trial_queue_4_D", "D")

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
