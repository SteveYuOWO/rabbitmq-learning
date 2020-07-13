package main

import "rabbitmq-learning/rabbitmq"

func main() {
	mq := rabbitmq.NewRabbitMQRouting("routingtest", "k1")
	mq.ConsumeRouting()
}
