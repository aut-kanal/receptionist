package mq

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gitlab.com/kanalbot/receptionist/configuration"
)

var (
	conn    *amqp.Connection
	channel *amqp.Channel

	qMsgs         amqp.Queue
	qAcceptedMsgs amqp.Queue
	qRejectedMsgs amqp.Queue

	acceptedMsgs <-chan amqp.Delivery
	rejectedMsgs <-chan amqp.Delivery
)

func SubscribeAcceptedMsgs(callback func(amqp.Delivery)) {
	go func() {
		for msg := range acceptedMsgs {
			go callback(msg)
		}
	}()
}

func SubscribeRejectedMsgs(callback func(amqp.Delivery)) {
	go func() {
		for msg := range rejectedMsgs {
			go callback(msg)
		}
	}()
}

func PublishMsg(data *amqp.Publishing) error {
	return channel.Publish("", qMsgs.Name, false, false, *data)
}

func InitMessageQueue() {
	// Connection
	var err error
	conn, err = amqp.Dial(configuration.GetInstance().GetString("rabbit-mq.url"))
	if err != nil {
		logrus.WithError(err).Fatalln("can't connect to message queue")
	}

	// Channel
	channel, err = conn.Channel()
	if err != nil {
		logrus.WithError(err).Fatalln("can't create mq channel")
	}

	err = channel.ExchangeDeclare(
		configuration.GetInstance().GetString("rabbit-mq.accept-ex-name"), // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		logrus.WithError(err).Fatal("can't declare accept exchange")
	}

	err = channel.ExchangeDeclare(
		configuration.GetInstance().GetString("rabbit-mq.reject-ex-name"), // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		logrus.WithError(err).Fatal("can't declare reject exchange")
	}

	// Queue
	qMsgs, err = channel.QueueDeclare(
		configuration.GetInstance().GetString("rabbit-mq.msg-queue-name"), // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		logrus.WithError(err).Fatalln("can't create messages queue")
	}

	qAcceptedMsgs, err = channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		logrus.WithError(err).Fatalln("can't create accepted messages queue")
	}

	qRejectedMsgs, err = channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		logrus.WithError(err).Fatalln("can't create rejected messages queue")
	}

	err = channel.QueueBind(
		qAcceptedMsgs.Name, // queue name
		"",                 // routing key
		configuration.GetInstance().GetString("rabbit-mq.accept-ex-name"), // exchange
		false,
		nil,
	)
	if err != nil {
		logrus.WithError(err).Fatal("can't bind queue to exchange")
	}

	err = channel.QueueBind(
		qRejectedMsgs.Name, // queue name
		"",                 // routing key
		configuration.GetInstance().GetString("rabbit-mq.reject-ex-name"), // exchange
		false,
		nil,
	)
	if err != nil {
		logrus.WithError(err).Fatal("can't bind queue to exchange")
	}

	acceptedMsgs, err = channel.Consume(
		qAcceptedMsgs.Name, // queue
		"",                 // consumer
		true,               // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	if err != nil {
		logrus.WithError(err).Fatal("can't init accepted msg consumer")
	}

	rejectedMsgs, err = channel.Consume(
		qRejectedMsgs.Name, // queue
		"",                 // consumer
		true,               // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	if err != nil {
		logrus.WithError(err).Fatal("can't init rejected msg consumer")
	}

	logrus.Info("message queue initialized")
}

func Close() {
	conn.Close()
}
