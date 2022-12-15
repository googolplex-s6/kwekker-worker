package main

import (
	"context"
	"fmt"
	"github.com/googolplex-s6/kwekker-protobufs/v3/kwek"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"kwekker-worker/util"
	"log"
	"time"
)

const kwekGuid = "f9d30d37-63a8-44a9-b2c3-3a45eb0701bc"
const kwekText = "Edited foo bar"

func main() {
	config, err := util.LoadConfig()
	if err != nil {
		log.Fatalln("Unable to load configuration; is the .env file present and valid?", err)
	}

	conn, err := amqp.Dial(
		fmt.Sprintf(
			"amqp://%s:%s@%s:%d%s",
			config.RabbitMQ.Username,
			config.RabbitMQ.Password,
			config.RabbitMQ.Host,
			config.RabbitMQ.Port,
			config.RabbitMQ.Vhost,
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

	err = ch.ExchangeDeclare("kweks", "topic", true, false, false, false, nil)
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

	err = ch.QueueBind(q.Name, "kwek.update", "kweks", false, nil)
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
		"kweks",
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
