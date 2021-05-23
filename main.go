package main

import "google-cloud-task-processor/amqp"

func main() {
	messageChannel := amqp.NewChannel()
	consumer := amqp.Consumer{MessageChannel: messageChannel}
	consumer.Consume()
}