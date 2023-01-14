package main

import (
	"context"
	"fmt"
	"github.com/googolplex-s6/kwekker-protobufs/v3/kwek"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"kwekker-worker/pkg/config"
	"log"
	"time"
)

const kwekGuid = "12230daf-29ee-47e0-b957-905e7731e12a"
const kwekText = "Edited foo bar"

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Unable to load configuration; is the .env file present and valid?", err)
	}

	conn, err := amqp.Dial(
		fmt.Sprintf(
			"amqp://%s:%s@%s:%d%s",
			conf.RabbitMQ.Username,
			conf.RabbitMQ.Password,
			conf.RabbitMQ.Host,
			conf.RabbitMQ.Port,
			conf.RabbitMQ.Vhost,
		),
	)

	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ", err)
	}

	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		log.Fatal("Failed to open channel", err)
	}

	defer ch.Close()

	err = ch.ExchangeDeclare("kwek-exchange", "topic", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare exchange", err)
	}

	updateKwekQueue(ch)
}

func updateKwekQueue(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(
		"kwek.update",
		true,
		false,
		false,
		false,
		nil,
	)

	err = ch.QueueBind(q.Name, "kwek.update", "kwek-exchange", false, nil)
	if err != nil {
		log.Fatal("Failed to bind queue", err)
	}

	if err != nil {
		log.Fatal("Failed to declare queue", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newKwek := &kwek.UpdateKwek{
		KwekGuid:  kwekGuid,
		Text:      kwekText,
		UpdatedAt: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
	}

	body, err := proto.Marshal(newKwek)
	if err != nil {
		log.Fatal("Failed to marshal protobuf", err)
	}

	err = ch.PublishWithContext(ctx,
		"kwek-exchange",
		q.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/protobuf",
			Body:         body,
		},
	)

	if err != nil {
		log.Fatal("Failed to publish message", err)
	}

	log.Println("Published message")
}
