package main

import "rabbitmq-learning/rabbitmq"

func main() {
	mq := rabbitmq.NewRabbitMQRouting("routingtest", "k2")
	mq.ConsumeRouting()
}
