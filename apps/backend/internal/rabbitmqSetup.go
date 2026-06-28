package internal

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Dipesh1203/alive/apps/backend/internal/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SetupRabbitMq() {
	rabbitmqURL := utils.GetEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672")
	log.Printf("[RABBITMQ] Attempting to connect to RabbitMQ at %s)", rabbitmqURL)
	conn, err2 := amqp.Dial(rabbitmqURL)
	log.Printf("[RABBITMQ] Connection result: %v", conn)

	utils.FailOnError(err2, "Fail to connect Rabbit MQ channel")
	log.Printf("[RABBITMQ] Opening channel...")

	ch, err3 := conn.Channel()
	utils.FailOnError(err3, "Fail to open channel")
	defer ch.Close()
	defer conn.Close()

	// log.Printf("[RABBITMQ] Declaring queue...")
	q, err4 := ch.QueueDeclare("Queue", false, false, false, false, nil)
	// log.Printf("[RABBITMQ] Queue declared: %v\", q)
	utils.FailOnError(err4, "Fail to declare Queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	body := "This my first Queue Entry"

	log.Printf("[RABBITMQ] Publishing message to queue...")
	err := ch.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{ContentType: "text/Plain",
		Body: []byte(body)})

	utils.FailOnError(err, "Fail o publish message")
	log.Printf("[RABBITMQ] Message published successfully: [X] Sent %s\n", body)

	log.Printf("[RABBITMQ] Registering consumer...")
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
	log.Printf("[RABBITMQ] Consumer registered, listening for messages...")

	for msg := range msgs {
		// log.Printf("[RABBITMQ] Received message body: %s\", string(msg.Body))
		fmt.Println("Body: ", string(msg.Body))
	}
	fmt.Printf("%T", msgs)
}
