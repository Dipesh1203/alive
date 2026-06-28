package services

import (
	"context"

	"github.com/Dipesh1203/alive/apps/backend/internal/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQService struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Queue   amqp.Queue
}

func NewRabbitMQService() (*RabbitMQService, error) {

	rabbitmqURL := utils.GetEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672")
	conn, err2 := amqp.Dial(rabbitmqURL)

	utils.FailOnError(err2, "Fail to connect Rabbit MQ channel")

	ch, err3 := conn.Channel()
	utils.FailOnError(err3, "Fail to open channel")

	q, err4 := ch.QueueDeclare("Queue", false, false, false, false, nil)
	// log.Printf("[RABBITMQ] Queue declared: %v\", q)
	utils.FailOnError(err4, "Fail to declare Queue")

	return &RabbitMQService{
		Conn:    conn,
		Channel: ch,
		Queue:   q,
	}, nil
}

func (r *RabbitMQService) PublishMessage(ctx context.Context, body []byte) error {
	return r.Channel.PublishWithContext(ctx, "", r.Queue.Name, false, false, amqp.Publishing{DeliveryMode: amqp.Persistent, ContentType: "text/Plain", Body: body})
}

func (r *RabbitMQService) Close() {
	if r.Channel != nil {
		r.Channel.Close()
	}
	if r.Conn != nil {
		r.Conn.Close()
	}
}
