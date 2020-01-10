package rabbitmq

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	channel  *amqp.Channel
	Name     string
	exchange string
}

func New(s string) *RabbitMQ {
	conn, err := amqp.Dial(s)
	if err != nil {
		logrus.Println("dial amqp error", err.Error())
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		logrus.Println("get amqp conn channel error", err.Error())
		panic(err)
	}

	q, err := ch.QueueDeclare(
		"",    //name 匿名的临时队列
		false, //durable
		true,  //delete when unused
		false, //exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		logrus.Println("channel declare queue failed", err.Error())
		panic(err)
	}

	mq := new(RabbitMQ)
	mq.channel = ch
	mq.Name = q.Name
	return mq
}

func (q *RabbitMQ) Bind(exchange string) {
	err := q.channel.QueueBind(
		q.Name,   //queue name
		"",       //routing key
		exchange, //exchange
		false,
		nil,
	)

	if err != nil {
		logrus.Println("queue bind queue failed", err.Error())
		panic(err)
	}

	q.exchange = exchange
}

func (q *RabbitMQ) Send(queue string, body interface{}) {
	str, err := json.Marshal(body)
	if err != nil {
		logrus.Println("mq send message json marshal body failed.", err.Error())
		panic(err)
	}

	err = q.channel.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ReplyTo: q.Name,
			Body:    []byte(str),
		})

	if err != nil {
		logrus.Println("amqp send message to queue single failed", err.Error())
		panic(err)
	}
}

func (q *RabbitMQ) Publish(exchange string, body interface{}) {
	str, err := json.Marshal(body)
	if err != nil {
		logrus.Println("mq publish msg json marshal failed.", err.Error())
		panic(err)
	}

	err = q.channel.Publish(
		exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ReplyTo: q.Name,
			Body:    []byte(str),
		})

	if err != nil {
		logrus.Println("publish msg failed", err.Error())
		panic(err)
	}
}

func (q *RabbitMQ) Consume() <-chan amqp.Delivery {
	c, err := q.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		logrus.Println("consume msg failed", err.Error())
		panic(err)
	}

	return c
}

func (q *RabbitMQ) Close() {
	q.channel.Close()
}
