package amqp

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"google-cloud-task-processor/shared"
	"log"
	"os"
)

type Consumer struct {
	MessageChannel <-chan amqp.Delivery
}

func (c *Consumer) Consume() {
	stopChan := make(chan bool)

	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())
		for d := range c.MessageChannel {
			log.Printf("Received a message: %s", d.Body)

			task := &shared.Task{}

			err := json.Unmarshal(d.Body, task)

			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}

			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("Acknowledged message")
			}

		}
	}()
	<-stopChan
}
