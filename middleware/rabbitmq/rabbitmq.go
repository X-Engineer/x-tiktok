package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

const MQURL = "amqp://guest:guest@106.14.252.145:5672/"

var BaseRmq *RabbitMQ

// RabbitMQ 参考: https://www.rabbitmq.com/tutorials/tutorial-one-go.html
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// 连接信息
	MqUrl string
}

func InitRabbitMQ() {
	BaseRmq = &RabbitMQ{
		MqUrl: MQURL,
	}
	conn, err := amqp.Dial(BaseRmq.MqUrl)
	BaseRmq.failOnError(err, "Failed to connect to RabbitMQ")
	BaseRmq.conn = conn
	BaseRmq.channel, err = conn.Channel()
	BaseRmq.failOnError(err, "Failed to get channel")
}

func (r *RabbitMQ) failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
		log.Panicf("%s: %s", msg, err)
	}
}

func (r *RabbitMQ) destroy() {
	r.channel.Close()
	r.conn.Close()
}
