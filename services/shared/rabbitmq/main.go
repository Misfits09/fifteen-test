package rabbitmq

import (
	"os"

	"github.com/labstack/gommon/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectToRabbitMq() (*amqp.Channel, amqp.Queue) {
	mq_url := os.Getenv("RABBITMQ_URL")
	conn, err := amqp.Dial(mq_url)

	if err != nil {
		log.Panic("Could not connect to Rabbit MQ", mq_url)
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
