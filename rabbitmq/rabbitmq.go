package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

// INSTALL
// docker run -dit --name rabbitmq -e RABBITMQ_DEFAULT_USER=admin -e RABBITMQ_DEFAULT_PASS=admin -p 15672:15672 -p 5672:5672 rabbitmq:management
// rabbitmqctl add_vhost [vhost_name] #create vhost
// rabbitmqctl delete_vhost [vhost_name] #delete vhost
// rabbitmqctl list_vhosts # list
// 配置最大连接限制，0：表示不可用，-1：无限制
// rabbitmqctl set_vhost_limits -p vhost_name '{"max-connections": 256}'
// 配置队列最大数，-1：无限制
// rabbitmqctl set_vhost_limits -p vhost_name '{"max-queues": 1024}'
// 权限
// rabbitmqctl set_permissions -p stevehost admin ".*" ".*" ".*"
// URL
// format: amqp://account:password@ip:port/vhost
const Url = "amqp://admin:admin@127.0.0.1:5672/stevehost"

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// 队列名称
	QueueName string
	// 交换机
	Exchange string
	// key
	Key string
	// URL
	URL string
}

// init RabbitMQ
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ{
	return &RabbitMQ{
		URL: Url,
		QueueName: queueName,
		Exchange: exchange,
		Key: key,
	}
}

func NewRabbitMQSimple(queueName string) *RabbitMQ{
	rabbitmq := NewRabbitMQ(queueName, "", "")
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.URL)
	rabbitmq.failOnErr(err, "amqp dial wrong")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "create channel wrong")
	return rabbitmq
}

func (r *RabbitMQ) PublishSimple(message string) {
	// 1.申请队列，如果不存在自动创建，存在就跳过
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		// 持久化
		false,
		// 自动删除
		false,
		// 排他性
		false,
		// 阻塞
		false,
		// 额外属性
		nil,
	)
	r.failOnErr(err, "declare queue error")
	err = r.channel.Publish(
		r.Exchange,
		r.QueueName,
		// 如果为true，则会根据exchange和routekey规则，如果未找到，会回退
		false,
		// 如果为true，exchange没有绑定消费者，则会返回给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	r.failOnErr(err, "publish message error")
}

func (r *RabbitMQ) ConsumeSimple() {
	// 1.申请队列，如果不存在自动创建，存在就跳过
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		// 持久化
		false,
		// 自动删除
		false,
		// 排他性
		false,
		// 阻塞
		false,
		// 额外属性
		nil,
	)
	r.failOnErr(err, "declare queue error")

	// consume
	consume, err := r.channel.Consume(
		r.QueueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	r.failOnErr(err, "consume message error")

	forever := make(chan bool)
	go func() {
		fmt.Println(consume)
		for d := range consume {
			log.Printf("Received: %s", d.Body)
		}
	}()
	log.Printf("[*] Waiting for message. To exit press CTRL*C")
	<- forever
}

// subscribe mode
func NewRabbitMQPubSub(exchange string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", exchange, "")
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.URL)
	rabbitmq.failOnErr(err, "amqp dial wrong")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "create channel wrong")
	return rabbitmq
}

// subscribe mode publish
func (r *RabbitMQ) PublishSubscribe(message string) {
	// 1.申请队列，如果不存在自动创建，存在就跳过
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		// 广播类型
		"fanout",
		// 持久化
		true,
		false,
		false,
		false,
		nil,
	)

	r.failOnErr(err, "declare queue error")

	err = r.channel.Publish(
		r.Exchange,
		"",
		// 如果为true，则会根据exchange和routekey规则，如果未找到，会回退
		false,
		// 如果为true，exchange没有绑定消费者，则会返回给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	r.failOnErr(err, "publish message error")
}

// subscribe mode consume
func (r *RabbitMQ) ConsumeSub() {
	// create exchange
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		// 广播类型
		"fanout",
		// 持久化
		true,
		false,
		false,
		false,
		nil,
	)

	r.failOnErr(err, "declare queue error")

	// create the queue
	q, err := r.channel.QueueDeclare(
		"",
		false,
		false,
		// 排他
		true,
		false,
		nil,
	)

	r.failOnErr(err, "declare queue error")

	// bind the queue
	err = r.channel.QueueBind(
		q.Name,
		"",
		r.Exchange,
		false,
		nil,
	)

	r.failOnErr(err, "consume message error")

	// consume
	consume, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	r.failOnErr(err, "consume message error")

	forever := make(chan bool)
	go func() {
		fmt.Println(consume)
		for d := range consume {
			log.Printf("Received: %s", d.Body)
		}
	}()
	log.Printf("[*] Waiting for message. To exit press CTRL*C")
	<- forever
}

// subscribe mode
func NewRabbitMQRouting(exchange string, routingKey string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", exchange, routingKey)
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.URL)
	rabbitmq.failOnErr(err, "amqp dial wrong")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "create channel wrong")
	return rabbitmq
}

// routing mode publish
func (r *RabbitMQ) PublishRouting(message string) {
	// 1.申请队列，如果不存在自动创建，存在就跳过
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		// 广播类型
		"direct",
		// 持久化
		true,
		false,
		false,
		false,
		nil,
	)

	r.failOnErr(err, "declare queue error")

	err = r.channel.Publish(
		r.Exchange,
		r.Key,
		// 如果为true，则会根据exchange和routekey规则，如果未找到，会回退
		false,
		// 如果为true，exchange没有绑定消费者，则会返回给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	r.failOnErr(err, "publish message error")
}


// routing mode consume
func (r *RabbitMQ) ConsumeRouting() {
	// create exchange
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		// 广播类型
		"direct",
		// 持久化
		true,
		false,
		false,
		false,
		nil,
	)

	r.failOnErr(err, "declare queue error")

	// create the queue
	q, err := r.channel.QueueDeclare(
		"",
		false,
		false,
		// 排他
		true,
		false,
		nil,
	)

	r.failOnErr(err, "declare queue error")

	// bind the queue
	err = r.channel.QueueBind(
		q.Name,
		r.Key,
		r.Exchange,
		false,
		nil,
	)

	r.failOnErr(err, "consume message error")

	// consume
	consume, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	r.failOnErr(err, "consume message error")

	forever := make(chan bool)
	go func() {
		fmt.Println(consume)
		for d := range consume {
			log.Printf("Received: %s", d.Body)
		}
	}()
	log.Printf("[*] Waiting for message. To exit press CTRL*C")
	<- forever
}

// destroy channel
func (r *RabbitMQ) Destroy() {
	r.channel.Close()
}

// handle error
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("message: %s, error: %s", message, err)
		panic(fmt.Sprintf("message: %s, error: %s", message, err))
	}
}

