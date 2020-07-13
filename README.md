# ⛵︎Rabbitmq-learning

### 1.⛪︎Intro

一个基于rabbitmq的简单案例

### 2.⛳️Category

- rabbitmq：对基础rabbitmq封装数据结构和函数
- simplemode：简单模式，即1生产1消费模式
- workmode：工作模式，1生产仅可以1个客户端消费，开启多个客户端同时消费（进行负载）
- subscribemode：订阅模式，1生产可以多个客户端消费，开启多个客户端消费模式
- routingmode：除了指定exchangeName，还需要指定key进行分发，相当于多了一个交换可以指定不同的key进行分发

