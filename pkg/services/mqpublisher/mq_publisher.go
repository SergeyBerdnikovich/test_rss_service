package mqpublisher

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

type mqPublisher struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

// NewMQPublisher ...
func NewMQPublisher() *mqPublisher {
	var conn *amqp.Connection
	var ch *amqp.Channel

	for {
		var err error

		conn, err = amqp.Dial(rabbitMQURL())
		if err != nil {
			fmt.Println("Failed to connect to RabbitMQ", err)
			time.Sleep(10 * time.Second)
			continue
		}

		ch, err = conn.Channel()
		if err != nil {
			fmt.Println("Failed to open a channel", err)
			time.Sleep(10 * time.Second)
			continue
		}

		if err == nil {
			break
		}
	}

	return &mqPublisher{
		connection: conn,
		channel:    ch,
	}
}

// MQPublisher ...
var MQPublisher = NewMQPublisher()

func (p *mqPublisher) InitializeQueues() {
	_, err := p.channel.QueueDeclare(
		"rss.feed_urls", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare a queue")

	_, err = p.channel.QueueDeclare(
		"rss.feed_items", // name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	failOnError(err, "Failed to declare a queue")
}

func (p *mqPublisher) Publish(queueName string, message []byte) {
	err := p.channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		})
	failOnError(err, "Failed to publish a message")
}

func (p *mqPublisher) Consume(queueName string) <-chan amqp.Delivery {
	msgs, err := p.channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	failOnError(err, "Failed to register a consumer")

	return msgs
}

func (p *mqPublisher) Close() {
	p.channel.Close()
	p.connection.Close()
}

func rabbitMQURL() string {
	url := os.Getenv("RABBIT_MQ_URL")
	if len(url) == 0 {
		return "amqp://guest:guest@localhost:5672/"
	}

	return url
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
