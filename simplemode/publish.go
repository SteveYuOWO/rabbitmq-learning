package main

import (
	"fmt"
	"rabbitmq-learning/rabbitmq"
)

func main() {
	mq := rabbitmq.NewRabbitMQSimple("test")
	mq.PublishSimple("hello, steve yu")
	fmt.Println("yes")
}
