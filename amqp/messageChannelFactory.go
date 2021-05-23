package amqp

import (
	"fmt"
	"github.com/streadway/amqp"
	"google-cloud-task-processor/utilities"
)

func NewChannel() <-chan amqp.Delivery {
	conn, err := amqp.Dial(Config.AMQPConnectionURL)
	utilities.HandleError(err, "Can't connect to AMQP")

	amqpChannel, err := conn.Channel()
	utilities.HandleError(err, "Can't create an amqpChannel")

	queue, err := amqpChannel.QueueDeclare(Config.AMQPQueueName, true, false, false, false, nil)
	utilities.HandleError(err, fmt.Sprintf("Could not declare %s queue", Config.AMQPQueueName))

	err = amqpChannel.Qos(1, 0, false)
	utilities.HandleError(err, "Could not configure QoS")

	messageChannel, err := amqpChannel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	utilities.HandleError(err, "Could not register consumer")

	return messageChannel
}
