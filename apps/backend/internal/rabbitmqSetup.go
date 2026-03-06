package internal

import (
	"backend/internal/utils"
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SetupRabbitMq() {
	conn, err2 := amqp.Dial("amqp://guest:guest@localhost:5672")
	log.Println("Conn", conn)
	utils.FailOnError(err2, "Fail to connect Rabbit MQ channel")
	ch, err3 := conn.Channel()
	utils.FailOnError(err3, "Fail to open channel")
	defer ch.Close()
	defer conn.Close()

	q, err4 := ch.QueueDeclare("Queue", false, false, false, false, nil)
	log.Println("Q ", q)
	utils.FailOnError(err4, "Fail to declare Queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	body := "This my first Queue Entry"
	err := ch.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{ContentType: "text/Plain",
		Body: []byte(body)})

	utils.FailOnError(err, "Fail o publish message")
	log.Printf(" [X] Sent %s\n", body)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	utils.FailOnError(err, "Failed to register a consumer")
	log.Println("Recieve Message ", msgs)
	for msg := range msgs {
		fmt.Println("Body: ", string(msg.Body))
	}
	fmt.Printf("%T", msgs)
}
