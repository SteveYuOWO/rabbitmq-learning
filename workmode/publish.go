package main

import (
	"fmt"
	"rabbitmq-learning/rabbitmq"
	"strconv"
	"time"
)

func main() {
	mq := rabbitmq.NewRabbitMQSimple("test")
	for i := 0; i <= 100; i++ {
		mq.PublishSimple("Hello! " + strconv.Itoa(i))
		time.Sleep(time.Second * 1)
		fmt.Println(i)
	}
}
