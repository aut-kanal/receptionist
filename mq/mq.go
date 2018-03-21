package mq

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gitlab.com/kanalbot/receptionist/configuration"
)

var (
	conn    *amqp.Connection
	channel *amqp.Channel

	qMsgs amqp.Queue
)

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

	logrus.Info("message queue initialized")
}

func Close() {
	conn.Close()
}
