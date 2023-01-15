package rabbitmq

import (
	"fmt"
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

func (w *RabbitMQWorker) ListenToQueues(queues config.Queues, protobufchan chan<- proto.Message) {
	conn, err := w.connect()

	if err != nil {
		w.logger.Fatal("Failed to connect to RabbitMQ", zap.Error(err))
	}

	w.logger.Debug("Connected to RabbitMQ")

	defer conn.Close()

	mqchannel, err := conn.Channel()

	if err != nil {
		w.logger.Fatal("Failed to open channel", zap.Error(err))
	}

	defer mqchannel.Close()

	exchanges := extractExchanges(queues)

	w.declareExchanges(exchanges, mqchannel)
	w.declareAndBindQueues(queues, mqchannel)
	w.consumeQueues(queues, mqchannel, protobufchan)

	select {}
}

func (w *RabbitMQWorker) connect() (*amqp.Connection, error) {
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
		return nil, fmt.Errorf("failed to connect to RabbitMQ")
	}

	return conn, nil
}

func extractExchanges(queues config.Queues) []string {
	exchangeMap := make(map[string]bool)

	for _, queueData := range queues {
		exchangeMap[queueData.Exchange] = true
	}

	exchanges := make([]string, 0, len(exchangeMap))

	for exchange := range exchangeMap {
		exchanges = append(exchanges, exchange)
	}

	return exchanges
}

func (w *RabbitMQWorker) declareExchanges(exchanges []string, mqchannel *amqp.Channel) {
	for _, exchange := range exchanges {
		err := mqchannel.ExchangeDeclare(
			exchange,
			"topic",
			true,
			false,
			false,
			false,
			nil,
		)

		if err != nil {
			w.logger.Fatal("Failed to declare exchange", err)
		}
	}
}

func (w *RabbitMQWorker) declareAndBindQueues(queues config.Queues, mqchannel *amqp.Channel) {
	for queue, queueData := range queues {
		_, err := mqchannel.QueueDeclare(
			queue,
			true,
			false,
			false,
			false,
			nil,
		)

		if err != nil {
			w.logger.Fatal("Failed to declare queue", zap.Error(err))
		}

		err = mqchannel.QueueBind(
			queue,
			queue,
			queueData.Exchange,
			false,
			nil,
		)

		if err != nil {
			w.logger.Fatal("Failed to bind queue", zap.Error(err))
		}
	}
}

func (w *RabbitMQWorker) consumeQueues(queues config.Queues, mqchannel *amqp.Channel, protobufchan chan<- proto.Message) {
	for queue, queueData := range queues {
		msgs, err := mqchannel.Consume(
			queue,
			"",
			false,
			false,
			false,
			false,
			nil,
		)

		if err != nil {
			w.logger.Fatal("Failed to consume queue", zap.Error(err))
		}

		go w.handleMessages(msgs, queueData.Type, protobufchan)
	}
}

func (w *RabbitMQWorker) handleMessages(msgs <-chan amqp.Delivery, prototype proto.Message, protobufchan chan<- proto.Message) {
	for msg := range msgs {
		protobuf := proto.Clone(prototype)
		err := proto.Unmarshal(msg.Body, protobuf)

		if err != nil {
			w.logger.Error("Failed to unmarshal protobuf", zap.Error(err))
			_ = msg.Nack(false, false)
			continue
		}

		valid := validation.Validate(protobuf)

		if !valid.Valid {
			w.logger.Error("Failed to validate protobuf", zap.Strings("errors", valid.Errors))
			_ = msg.Nack(false, false)
			continue
		}

		protobufchan <- protobuf

		err = msg.Ack(true)

		if err != nil {
			w.logger.Error("Failed to acknowledge message", zap.Error(err))
		}
	}
}
