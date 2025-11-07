package queue

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQClient, yayıncı ve tüketici kanallarını yönetir.
type RabbitMQClient struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

// Connect, RabbitMQ sunucusuna bağlanır ve bir client nesnesi döner.
func Connect(url string) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq connection failed: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	log.Println("✓ RabbitMQ connected and channel opened")
	return &RabbitMQClient{Conn: conn, Channel: ch}, nil
}

// Close, kanalı ve bağlantıyı düzgün bir şekilde kapatır.
func (c *RabbitMQClient) Close() {
	if c.Channel != nil {
		c.Channel.Close()
	}
	if c.Conn != nil {
		c.Conn.Close()
	}
}

// DeclareAndBindQueue, bir exchange, queue oluşturur ve bunları birbirine bağlar.
// Bu fonksiyon "idempotent"tir, yani kaynaklar zaten varsa hata vermeden devam eder.
func (c *RabbitMQClient) DeclareAndBindQueue(exchange, queueName, routingKey string) error {
	// Dayanıklı (durable) bir exchange tanımlıyoruz. Sunucu yeniden başlasa bile kaybolmaz.
	err := c.Channel.ExchangeDeclare(
		exchange, // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return fmt.Errorf("exchange declaration failed: %w", err)
	}

	// Dayanıklı bir kuyruk tanımlıyoruz.
	q, err := c.Channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("queue declaration failed: %w", err)
	}

	// Kuyruğu exchange'e routing key ile bağlıyoruz.
	err = c.Channel.QueueBind(
		q.Name,     // queue name
		routingKey, // routing key
		exchange,   // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue bind failed: %w", err)
	}

	log.Printf("✓ Queue '%s' declared and bound to exchange '%s'", queueName, exchange)
	return nil
}

// Publish, bir mesajı belirtilen exchange'e ve routing key'e gönderir.
func (c *RabbitMQClient) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := c.Channel.PublishWithContext(ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Mesajın diske yazılmasını sağlar
		})

	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	return nil
}
