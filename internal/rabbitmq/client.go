package rabbitmq

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Vhost    string
}

func (c *Config) URL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		c.Username, c.Password, c.Host, c.Port, c.Vhost)
}

func NewClient(cfg *Config) (*Client, error) {
	conn, err := amqp.Dial(cfg.URL())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &Client{
		conn:    conn,
		channel: ch,
	}, nil
}

func (c *Client) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) Publish(exchange, routingKey, body string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.channel.PublishWithContext(ctx,
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
			Timestamp:   time.Now(),
		},
	)
}

func (c *Client) PublishToQueue(queueName, body string) error {
	_, err := c.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	return c.Publish("", queueName, body)
}

type Message struct {
	// Method Frame info
	ConsumerTag string
	DeliveryTag uint64
	Redelivered bool
	Exchange    string
	RoutingKey  string

	// Header Frame info
	ContentType     string
	ContentEncoding string
	Headers         map[string]interface{}
	DeliveryMode    uint8
	Priority        uint8
	CorrelationId   string
	ReplyTo         string
	Expiration      string
	MessageId       string
	Timestamp       time.Time
	Type            string
	UserId          string
	AppId           string

	// Body Frame info
	Body    string
	RawBody []byte
}

func (c *Client) Consume(queueName string, autoAck bool, handler func(msg Message) error) error {
	_, err := c.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	msgs, err := c.channel.Consume(
		queueName,
		"",
		autoAck,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	for d := range msgs {
		msg := Message{
			// Method Frame
			ConsumerTag: d.ConsumerTag,
			DeliveryTag: d.DeliveryTag,
			Redelivered: d.Redelivered,
			Exchange:    d.Exchange,
			RoutingKey:  d.RoutingKey,

			// Header Frame
			ContentType:     d.ContentType,
			ContentEncoding: d.ContentEncoding,
			Headers:         d.Headers,
			DeliveryMode:    d.DeliveryMode,
			Priority:        d.Priority,
			CorrelationId:   d.CorrelationId,
			ReplyTo:         d.ReplyTo,
			Expiration:      d.Expiration,
			MessageId:       d.MessageId,
			Timestamp:       d.Timestamp,
			Type:            d.Type,
			UserId:          d.UserId,
			AppId:           d.AppId,

			// Body Frame
			Body:    string(d.Body),
			RawBody: d.Body,
		}

		if err := handler(msg); err != nil {
			if !autoAck {
				d.Nack(false, true)
			}
			continue
		}

		if !autoAck {
			d.Ack(false)
		}
	}

	return nil
}
