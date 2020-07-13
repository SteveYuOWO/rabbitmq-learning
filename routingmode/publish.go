package main

import (
	"fmt"
	"rabbitmq-learning/rabbitmq"
	"strconv"
	"time"
)

func main() {
	k1 := rabbitmq.NewRabbitMQRouting("routingtest", "k1")
	k2 := rabbitmq.NewRabbitMQRouting("routingtest", "k2")
	for i := 0; i <= 100; i++ {
		k1.PublishRouting("Hello! " + strconv.Itoa(i))
		k2.PublishRouting("Hello! " + strconv.Itoa(i))
		time.Sleep(time.Second * 1)
		fmt.Println(i)
	}
}
