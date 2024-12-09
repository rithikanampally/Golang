package queue

import (
	"github.com/streadway/amqp"
)

// InitRabbitMQ initializes a RabbitMQ connection and channel.
func InitRabbitMQ(connString string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(connString)
	if err != nil {
		return nil, nil, err
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}
	return conn, channel, nil
}

// PublishToQueue sends a message to the specified RabbitMQ queue.
func PublishToQueue(channel *amqp.Channel, queueName string, message []byte) error {
	_, err := channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}
	return channel.Publish("", queueName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        message,
	})
}
