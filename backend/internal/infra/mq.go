package infra

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"bigmarket/internal/pkg"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Event struct {
	ID        string      `json:"id"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

type MQ struct {
	url  string
	mu   sync.Mutex
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewMQ(url string) *MQ { return &MQ{url: url} }

func (m *MQ) ensure() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.ch != nil && !m.ch.IsClosed() {
		return nil
	}
	var err error
	for i := 0; i < 30; i++ {
		m.conn, err = amqp.Dial(m.url)
		if err == nil {
			m.ch, err = m.conn.Channel()
			if err == nil {
				err = m.ch.ExchangeDeclare(pkg.ExchangeName, "topic", true, false, false, false, nil)
				if err == nil {
					return nil
				}
			}
		}
		slog.Warn("waiting for rabbitmq", "err", err)
		time.Sleep(2 * time.Second)
	}
	return err
}

func (m *MQ) Publish(topic string, data interface{}) error {
	if err := m.ensure(); err != nil {
		return err
	}
	body, err := json.Marshal(Event{ID: pkg.GenMessageID(), Timestamp: time.Now().UnixMilli(), Data: data})
	if err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.ch.Publish(pkg.ExchangeName, topic, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
	})
}

// PublishRaw 发布已序列化的事件消息（任务表补偿时使用）
func (m *MQ) PublishRaw(topic, message string) error {
	if err := m.ensure(); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.ch.Publish(pkg.ExchangeName, topic, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         []byte(message),
	})
}

func (m *MQ) Subscribe(topic string, handler func([]byte) error) {
	for {
		conn, err := amqp.Dial(m.url)
		if err != nil {
			slog.Error("mq subscribe connect failed", "topic", topic, "err", err)
			time.Sleep(5 * time.Second)
			continue
		}
		ch, err := conn.Channel()
		if err == nil {
			err = ch.ExchangeDeclare(pkg.ExchangeName, "topic", true, false, false, false, nil)
		}
		if err == nil {
			_, err = ch.QueueDeclare(topic, true, false, false, false, nil)
		}
		if err == nil {
			err = ch.QueueBind(topic, topic, pkg.ExchangeName, false, nil)
		}
		if err != nil {
			slog.Error("mq declare queue failed", "topic", topic, "err", err)
			_ = conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}
		msgs, err := ch.Consume(topic, "", false, false, false, false, nil)
		if err != nil {
			slog.Error("mq consume failed", "topic", topic, "err", err)
			_ = conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}
		slog.Info("mq consumer started", "topic", topic)
		for msg := range msgs {
			if err := handler(msg.Body); err != nil {
				slog.Error("mq handle failed", "topic", topic, "err", err)
				_ = msg.Nack(false, true)
				continue
			}
			_ = msg.Ack(false)
		}
		slog.Warn("mq channel closed, reconnecting", "topic", topic)
		_ = conn.Close()
		time.Sleep(3 * time.Second)
	}
}

func (m *MQ) String() string { return fmt.Sprintf("MQ(%s)", m.url) }
