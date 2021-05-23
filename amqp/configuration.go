package amqp

type Configuration struct {
	AMQPConnectionURL string
	AMQPQueueName string
}

var Config = Configuration{
	AMQPConnectionURL: "amqp://guest:guest@localhost:5672/",
	AMQPQueueName: "queue",
}
