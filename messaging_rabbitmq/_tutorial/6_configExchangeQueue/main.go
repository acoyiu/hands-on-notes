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

	// Exchange Config

	// Features in MQ Manager are meaning these:
	// 1: Durable     + Non-Auto-Deleted 	  - Always Not Kill
	// 2. Durable     + Auto-Deleted		  - Kill when no bind at server start
	// 3. Non-Durable + Non-Auto-deleted      - Kill on restart Only
	// 4: Non-Durable + Auto-Deleted 	      - Kill on restart OR Kill on No Binding

	// 'amq.fanout' default is durable

	err = mqChannel.ExchangeDeclare(
		"trial_exchange_6",  // name of ex
		amqp.ExchangeFanout, // type 		  - way to distribute message
		true,                // durable
		false,               // auto-deleted
		false,               // internal 	  - can not published fron outside
		false,               // no-wait       - declare without waiting confirmation from the server, may cause error
		nil,                 // arguments     - amqp.Table for Extra parameters
	)
	if err != nil {
		fatalError(err, "Unable to create Exchange in RabbitMQ error")
	}

	// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

	// Queue Config

	// 1: exclusive - 聲明該 queue 的 connection 才可 access，其他 connection 所有 action 都會 error
	// 2: noWait    - 假定 queue 已存在於 server，不會 redeclare，若無，會 error
	// 3: args      - Extra Configs

	mqChannel.QueueDeclare(
		"trial_queue_6", // name string
		false,           // durable bool
		true,            // autoDelete bool
		true,            // exclusive bool
		false,           // noWait bool
		nil,             // args Table
	)
}
