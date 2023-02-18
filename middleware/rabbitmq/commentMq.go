package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"strconv"
	"x-tiktok/dao"
)

type CommentMQ struct {
	RabbitMQ
	// 队列的名称
	QueueName string
	// 交换机的名称
	Exchange string
	// bind key 的名称
	Key string
}

var SimpleCommentDelMQ *CommentMQ

func (r *CommentMQ) PublishSimple(message string) error {
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

// ConsumeSimple simple 模式下消费者
func (r *CommentMQ) ConsumeSimple() {
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
	// 由于目前消息只能传递string类型，插入评论操作不适合写入消息队列
	go r.consumerCommentDel(msgs)
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// 添加删除评论的消费
func (r *CommentMQ) consumerCommentDel(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		// 解析参数
		cId := fmt.Sprintf("%s", msg.Body)
		log.Println("添加评论删除消费者获得 cId:", cId)
		commentId, _ := strconv.ParseInt(cId, 10, 64)
		// 数据库操作，最大重试次数 cnt
		cnt := 10
		for i := 0; i < cnt; i++ {
			succeed := true
			var err error
			err = dao.DeleteComment(commentId)
			if err != nil {
				succeed = false
			}
			if succeed {
				break
			}
		}
	}

}

// 新建评论消息队列
func newCommentRabbitMQ(queueName string, exchangeName string, key string) *CommentMQ {
	commentMq := &CommentMQ{
		RabbitMQ:  *BaseRmq,
		QueueName: queueName,
		Exchange:  exchangeName,
		Key:       key,
	}
	return commentMq
}

// NewSimpleCommentRabbitMQ 新建简单模式的消息队列（生产者，消息队列，一个消费者）
func NewSimpleCommentRabbitMQ(queueName string) *CommentMQ {
	return newCommentRabbitMQ(queueName, "", "")
}

func InitCommentRabbitMQ() {
	SimpleCommentDelMQ = NewSimpleCommentRabbitMQ("comment_del")
	// 开启 go routine 启动消费者
	go SimpleCommentDelMQ.ConsumeSimple()
}
