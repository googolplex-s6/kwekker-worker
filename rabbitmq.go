package main

import (
	"fmt"
	kwekker_protobufs "github.com/googolplex-s6/kwekker-protobufs/v2/kwek"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"kwekker-worker/util"
)

type RabbitMQWorker struct {
	logger *zap.SugaredLogger
	config *util.RabbitMQConfig
}

func NewRabbitMQWorker(logger *zap.SugaredLogger, config *util.RabbitMQConfig) *RabbitMQWorker {
	return &RabbitMQWorker{
		logger: logger,
		config: config,
	}
}

func (w *RabbitMQWorker) Listen(
	createKwekChannel chan<- *kwekker_protobufs.CreateKwek,
	updateKwekChannel chan<- *kwekker_protobufs.UpdateKwek,
	deleteKwekChannel chan<- *kwekker_protobufs.DeleteKwek,
) {
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

	go w.createKwekQueue(createKwekChannel, ch)
	go w.updateKwekQueue(updateKwekChannel, ch)
	go w.deleteKwekQueue(deleteKwekChannel, ch)

	select {}
}

func (w *RabbitMQWorker) createKwekQueue(kwekChannel chan<- *kwekker_protobufs.CreateKwek, ch *amqp.Channel) {
	q, err := ch.QueueDeclare(
		"kwek.create",
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
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		w.logger.Fatal("Failed to register consumer", zap.Error(err))
	}

	w.logger.Info("Listening for kweks...")

	go func() {
		for d := range msgs {
			w.logger.Debug("Received message")

			kwek := &kwekker_protobufs.CreateKwek{}
			err := proto.Unmarshal(d.Body, kwek)

			if err != nil {
				w.logger.Error("Failed to unmarshal kwek", zap.Error(err))
				continue
			}

			kwekChannel <- kwek

			err = d.Ack(true)

			if err != nil {
				w.logger.Error("Failed to acknowledge message", zap.Error(err))
			}
		}
	}()

	select {}
}

func (w *RabbitMQWorker) updateKwekQueue(kwekChannel chan<- *kwekker_protobufs.UpdateKwek, ch *amqp.Channel) {
	q, err := ch.QueueDeclare(
		"kwek.update",
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
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		w.logger.Fatal("Failed to register consumer", zap.Error(err))
	}

	w.logger.Info("Listening for kweks...")

	go func() {
		for d := range msgs {
			w.logger.Debug("Received message")

			kwek := &kwekker_protobufs.UpdateKwek{}
			err := proto.Unmarshal(d.Body, kwek)

			if err != nil {
				w.logger.Error("Failed to unmarshal kwek", zap.Error(err))
				continue
			}

			kwekChannel <- kwek

			err = d.Ack(true)

			if err != nil {
				w.logger.Error("Failed to acknowledge message", zap.Error(err))
			}
		}
	}()

	select {}
}

func (w *RabbitMQWorker) deleteKwekQueue(kwekChannel chan<- *kwekker_protobufs.DeleteKwek, ch *amqp.Channel) {
	q, err := ch.QueueDeclare(
		"kwek.delete",
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
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		w.logger.Fatal("Failed to register consumer", zap.Error(err))
	}

	w.logger.Info("Listening for kweks...")

	go func() {
		for d := range msgs {
			w.logger.Debug("Received message")

			kwek := &kwekker_protobufs.DeleteKwek{}
			err := proto.Unmarshal(d.Body, kwek)

			if err != nil {
				w.logger.Error("Failed to unmarshal kwek", zap.Error(err))
				continue
			}

			kwekChannel <- kwek

			err = d.Ack(true)

			if err != nil {
				w.logger.Error("Failed to acknowledge message", zap.Error(err))
			}
		}
	}()

	select {}
}
