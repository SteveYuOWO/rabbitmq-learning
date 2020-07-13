package main

import "rabbitmq-learning/rabbitmq"

func main() {
	mq := rabbitmq.NewRabbitMQSimple("test")
	mq.ConsumeSimple()
}
