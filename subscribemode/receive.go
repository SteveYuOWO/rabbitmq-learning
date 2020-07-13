package main

import "rabbitmq-learning/rabbitmq"

func main() {
	mq := rabbitmq.NewRabbitMQPubSub("pubtest")
	mq.ConsumeSub()
}
