package main

import (
	"fmt"
	kwekker_protobufs "github.com/googolplex-s6/kwekker-protobufs/kwek"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"kwekker-worker/util"
)

type RabbitMQWorker struct {
	logger *zap.Logger
	config *util.RabbitMQConfig
}

func NewRabbitMQWorker(logger *zap.Logger, config *util.RabbitMQConfig) *RabbitMQWorker {
	return &RabbitMQWorker{
		logger: logger,
		config: config,
	}
}

func (w *RabbitMQWorker) Listen(kwekChannel chan<- *kwekker_protobufs.Kwek) {
	conn, err := amqp.Dial(
		fmt.Sprintf(
			"amqp://%s:%s@%s:%d%s",
			w.config.Username,
			w.config.Password,
			w.config.Host,
			w.config.Port,
			w.config.Vhost,
		),
	)

	if err != nil {
		w.logger.Fatal("Failed to connect to RabbitMQ", zap.Error(err))
	}

	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		w.logger.Fatal("Failed to open channel", zap.Error(err))
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		"kweks",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		w.logger.Fatal("Failed to declare queue", zap.Error(err))
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true, // TODO: Ack manually later
		false,
		false,
		false,
		nil,
	)

	w.logger.Info("Listening for kweks...")

	go func() {
		for d := range msgs {
			w.logger.Debug("Received message")

			kwek := &kwekker_protobufs.Kwek{}
			err := proto.Unmarshal(d.Body, kwek)

			if err != nil {
				w.logger.Error("Failed to unmarshal kwek", zap.Error(err))
				continue
			}

			kwekChannel <- kwek
		}
	}()

	select {} // TODO: Handle shutdown
}
