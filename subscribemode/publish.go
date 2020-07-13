package main

import (
	"fmt"
	"rabbitmq-learning/rabbitmq"
	"strconv"
	"time"
)

func main() {
	mq := rabbitmq.NewRabbitMQPubSub("pubtest")
	for i := 0; i <= 100; i++ {
		mq.PublishSubscribe("Hello! " + strconv.Itoa(i))
		time.Sleep(time.Second * 1)
		fmt.Println(i)
	}
}
