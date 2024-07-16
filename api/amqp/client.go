package amqp

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log/slog"
	"os"

	"github.com/aacebo/orca/api/utils"

	"github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn   *amqp091.Connection
	ch     *amqp091.Channel
	log    *slog.Logger
	queues map[string]amqp091.Queue
}

func New() *Client {
	self := &Client{
		conn:   nil,
		ch:     nil,
		log:    slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("name", "amqp"),
		queues: map[string]amqp091.Queue{},
	}

	return self.Connect()
}

func (self *Client) Connect() *Client {
	if self.conn == nil || self.conn.IsClosed() {
		conn, err := amqp091.Dial(fmt.Sprintf(
			"amqp://%s:%s@%s:%s/",
			utils.GetEnv("RABBIT_USER", "admin"),
			utils.GetEnv("RABBIT_PASSWORD", "admin"),
			utils.GetEnv("RABBIT_HOST", "localhost"),
			utils.GetEnv("RABBIT_PORT", "5672"),
		))

		if err != nil {
			self.log.Error(err.Error())
		}

		self.conn = conn
		self.log.Info("connection established...")
	}

	if self.ch == nil || self.ch.IsClosed() {
		ch, err := self.conn.Channel()

		if err != nil {
			self.log.Error(err.Error())
		}

		self.ch = ch
	}

	return self
}

func (self *Client) Close() {
	self.ch.Close()
	self.conn.Close()
}

func (self Client) Closed() bool {
	if self.conn == nil || self.conn.IsClosed() ||
		self.ch == nil || self.ch.IsClosed() {
		return true
	}

	return false
}

func (self Client) Publish(exchange string, queue string, body any) {
	key := fmt.Sprintf("%s.%s", exchange, queue)

	if self.Closed() {
		self.Connect()
	}

	_, exists := self.queues[key]

	if !exists {
		q, err := self.assertQueue(exchange, queue)

		if err != nil {
			self.log.Error(err.Error())
		}

		self.queues[key] = *q
	}

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(body)

	if err != nil {
		self.log.Error(err.Error())
	}

	err = self.ch.PublishWithContext(
		context.Background(),
		exchange,
		key,
		false,
		false,
		amqp091.Publishing{
			DeliveryMode: amqp091.Persistent,
			Body:         buf.Bytes(),
		},
	)

	if err != nil {
		self.log.Error(err.Error())
	}
}

func (self *Client) Consume(exchange string, queue string, handler func(amqp091.Delivery)) {
	key := fmt.Sprintf("%s.%s", exchange, queue)

	if self.Closed() {
		self.Connect()
	}

	_, exists := self.queues[key]

	if !exists {
		q, err := self.assertQueue(exchange, queue)

		if err != nil {
			self.log.Error(err.Error())
		}

		self.queues[key] = *q
	}

	msgs, err := self.ch.Consume(
		key,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		self.log.Error(err.Error())
	}

	go func() {
		for {
			msg, ok := <-msgs

			if !ok || self.Closed() {
				self.Consume(exchange, queue, handler)
				break
			}

			go handler(msg)
		}
	}()
}

func (self Client) assertQueue(exchange string, queue string) (*amqp091.Queue, error) {
	key := fmt.Sprintf("%s.%s", exchange, queue)
	err := self.ch.ExchangeDeclare(
		exchange,
		"topic",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	q, err := self.ch.QueueDeclare(
		key,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	err = self.ch.QueueBind(
		key,
		key,
		exchange,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	return &q, nil
}
