package rabbitmq

import (
	"fmt"
	kwekkerprotobufs "github.com/googolplex-s6/kwekker-protobufs/v3/kwek"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"kwekker-worker/pkg/config"
	"kwekker-worker/pkg/validation"
	"time"
)

type RabbitMQWorker struct {
	logger *zap.SugaredLogger
	config config.RabbitMQConfig
}

func NewRabbitMQWorker(logger *zap.SugaredLogger, config config.RabbitMQConfig) *RabbitMQWorker {
	return &RabbitMQWorker{
		logger: logger,
		config: config,
	}
}

func (w *RabbitMQWorker) Listen(
	createKwekChannel chan<- *kwekkerprotobufs.CreateKwek,
	updateKwekChannel chan<- *kwekkerprotobufs.UpdateKwek,
	deleteKwekChannel chan<- *kwekkerprotobufs.DeleteKwek,
) {
	var conn *amqp.Connection

	for i := 0; i < 5; i++ {
		var err error
		conn, err = amqp.Dial(
			fmt.Sprintf(
				"amqp://%s:%s@%s:%d%s",
				w.config.Username,
				w.config.Password,
				w.config.Host,
				w.config.Port,
				w.config.Vhost,
			),
		)

		if err == nil {
			break
		}

		w.logger.Debug("Retrying RabbitMQ connection", zap.Error(err))
		time.Sleep(5 * time.Second)
	}

	if conn == nil {
		w.logger.Fatal("Failed to connect to RabbitMQ")
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

func (w *RabbitMQWorker) createKwekQueue(kwekChannel chan<- *kwekkerprotobufs.CreateKwek, ch *amqp.Channel) {
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
			w.handleCreateKwekDelivery(d, kwekChannel)
		}
	}()

	select {}
}

func (w *RabbitMQWorker) updateKwekQueue(kwekChannel chan<- *kwekkerprotobufs.UpdateKwek, ch *amqp.Channel) {
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
			w.handleUpdateKwekDelivery(d, kwekChannel)
		}
	}()

	select {}
}

func (w *RabbitMQWorker) deleteKwekQueue(kwekChannel chan<- *kwekkerprotobufs.DeleteKwek, ch *amqp.Channel) {
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
			w.handleDeleteKwekDelivery(d, kwekChannel)
		}
	}()

	select {}
}

func (w *RabbitMQWorker) handleCreateKwekDelivery(d amqp.Delivery, c chan<- *kwekkerprotobufs.CreateKwek) {
	kwek := &kwekkerprotobufs.CreateKwek{}
	err := proto.Unmarshal(d.Body, kwek)

	if err != nil {
		w.logger.Error("Failed to unmarshal kwek", zap.Error(err))
		_ = d.Nack(false, false)
		return
	}

	valid := validation.ValidateCreateKwek(kwek)

	if !valid.Valid {
		w.logger.Error("Failed to validate kwek", zap.Strings("errors", valid.Errors))
		_ = d.Nack(false, false)
		return
	}

	c <- kwek

	err = d.Ack(true)

	if err != nil {
		w.logger.Error("Failed to acknowledge message", zap.Error(err))
	}
}

func (w *RabbitMQWorker) handleUpdateKwekDelivery(d amqp.Delivery, c chan<- *kwekkerprotobufs.UpdateKwek) {
	kwek := &kwekkerprotobufs.UpdateKwek{}
	err := proto.Unmarshal(d.Body, kwek)

	if err != nil {
		w.logger.Error("Failed to unmarshal kwek", zap.Error(err))
		_ = d.Nack(false, false)
		return
	}

	valid := validation.ValidateUpdateKwek(kwek)

	if !valid.Valid {
		w.logger.Error("Failed to validate kwek", zap.Strings("errors", valid.Errors))
		_ = d.Nack(false, false)
		return
	}

	c <- kwek

	err = d.Ack(true)

	if err != nil {
		w.logger.Error("Failed to acknowledge message", zap.Error(err))
	}
}

func (w *RabbitMQWorker) handleDeleteKwekDelivery(d amqp.Delivery, c chan<- *kwekkerprotobufs.DeleteKwek) {
	kwek := &kwekkerprotobufs.DeleteKwek{}
	err := proto.Unmarshal(d.Body, kwek)

	if err != nil {
		w.logger.Error("Failed to unmarshal kwek", zap.Error(err))
		_ = d.Nack(false, false)
		return
	}

	valid := validation.ValidateDeleteKwek(kwek)

	if !valid.Valid {
		w.logger.Error("Failed to validate kwek", zap.Strings("errors", valid.Errors))
		_ = d.Nack(false, false)
		return
	}

	c <- kwek

	err = d.Ack(true)

	if err != nil {
		w.logger.Error("Failed to acknowledge message", zap.Error(err))
	}
}
