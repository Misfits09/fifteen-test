package rabbitmq

import (
	"os"

	"github.com/labstack/gommon/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectToRabbitMq() (*amqp.Channel, amqp.Queue) {
	mqURL := os.Getenv("RABBITMQ_URL")
	conn, err := amqp.Dial(mqURL)

	if err != nil {
		log.Panic("Could not connect to Rabbit MQ")
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Panic("Could not connect to Rabbit MQ Channel")
	}

	queue, err := channel.QueueDeclare("location_update", true, false, false, false, nil)
	if err != nil {
		log.Panic("Could declare queue")
	}

	return channel, queue
}
