package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"strconv"
	"strings"
	"x-tiktok/config"
	"x-tiktok/dao"
)

type FollowMQ struct {
	RabbitMQ
	// 队列的名称
	QueueName string
	// 交换机的名称
	Exchange string
	// bind key 的名称
	Key string
}

var SimpleFollowAddMQ *FollowMQ
var SimpleFollowDelMQ *FollowMQ

// PublishSimpleFollow simple 模式下关注生产者
func (r *FollowMQ) PublishSimpleFollow(message string) error {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞处理
		false,
		//额外的属性
		nil,
	)
	if err != nil {
		log.Println(err)
		return err
	}
	//调用channel 发送消息到队列中
	err = r.channel.Publish(
		r.Exchange,
		r.QueueName,
		//如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,
		//如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		return err
	}
	return nil
}

// ConsumeSimpleFollow simple 模式下消费者 follow模块
func (r *FollowMQ) ConsumeSimpleFollow() {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞处理
		false,
		//额外的属性
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	//接收消息
	msgs, err := r.channel.Consume(
		q.Name, // queue
		//用来区分多个消费者
		"", // consumer
		//是否自动应答
		true, // auto-ack
		//是否独有
		false, // exclusive
		//设置为true，表示 不能将同一个Connection中生产者发送的消息传递给这个Connection中 的消费者
		false, // no-local
		//列是否阻塞
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
	//启用协程处理消息
	//go func() {
	//	for d := range msgs {
	//		//消息逻辑处理，可以自行设计逻辑
	//		log.Printf("Received a message: %s", d.Body)
	//	}
	//}()

	log.Println("q.Name", q.Name)
	switch q.Name {
	case "follow_add":
		go r.consumerFollowAdd(msgs)
	case "follow_del":
		go r.consumerFollowDel(msgs)
	}

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// consumerFollowAdd 添加关注关系的消费
func (r *FollowMQ) consumerFollowAdd(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		followDao := dao.NewFollowDaoInstance()
		// 解析参数
		params := strings.Split(fmt.Sprintf("%s", msg.Body), "-")
		log.Println("添加关注关系消费者获得 params:", params)
		userId, _ := strconv.ParseInt(params[0], 10, 64)
		targetId, _ := strconv.ParseInt(params[1], 10, 64)
		insertOrUpdate := params[2]
		// 数据库操作，最大重试次数 cnt
		cnt := 10
		for i := 0; i < cnt; i++ {
			succeed := true
			var err error
			if insertOrUpdate == config.DB_INSERT {
				_, err = followDao.InsertFollowRelation(userId, targetId)
			} else if insertOrUpdate == config.DB_UPDATE {
				_, err = followDao.UpdateFollowRelation(userId, targetId, int8(1))
			}
			if err != nil {
				succeed = false
			}
			if succeed {
				break
			}
		}
	}
}

// consumerFollowDel 删除关注关系的消费
func (r *FollowMQ) consumerFollowDel(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		followDao := dao.NewFollowDaoInstance()
		// 解析参数
		params := strings.Split(fmt.Sprintf("%s", msg.Body), "-")
		log.Println("添加关注关系消费者获得 params:", params)
		userId, _ := strconv.ParseInt(params[0], 10, 64)
		targetId, _ := strconv.ParseInt(params[1], 10, 64)
		insertOrUpdate := params[2]
		// 数据库操作，最大重试次数 cnt
		cnt := 10
		for i := 0; i < cnt; i++ {
			succeed := true
			var err error
			if insertOrUpdate == config.DB_UPDATE {
				_, err = followDao.UpdateFollowRelation(userId, targetId, int8(0))
			}
			if err != nil {
				succeed = false
			}
			if succeed {
				break
			}
		}
	}
}

// 新建 "关注" 消息队列
func newFollowRabbitMQ(queueName string, exchangeName string, key string) *FollowMQ {
	followMQ := &FollowMQ{
		RabbitMQ:  *BaseRmq,
		QueueName: queueName,
		Exchange:  exchangeName,
		Key:       key,
	}
	return followMQ
}

// NewSimpleFollowRabbitMQ 新建简单模式的消息队列（生产者，消息队列，一个消费者）
func NewSimpleFollowRabbitMQ(queueName string) *FollowMQ {
	return newFollowRabbitMQ(queueName, "", "")
}

func InitFollowRabbitMQ() {
	SimpleFollowAddMQ = NewSimpleFollowRabbitMQ("follow_add")
	SimpleFollowDelMQ = NewSimpleFollowRabbitMQ("follow_del")
	// 开启 go routine 启动消费者
	go SimpleFollowAddMQ.ConsumeSimpleFollow()
	go SimpleFollowDelMQ.ConsumeSimpleFollow()
}
