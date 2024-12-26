package producer

import (
	"kafka-demo/internal/broker"
	"kafka-demo/pkg/types"
	"time"

	"github.com/google/uuid"
)

type Producer struct {
	broker *broker.Broker
	config types.Config
}

func NewProducer(broker *broker.Broker, config types.Config) *Producer {
	return &Producer{
		broker: broker,
		config: config,
	}
}

func (p *Producer) Send(topic string, payload []byte) error {
	msg := &types.Message{
		ID:        uuid.New().String(),
		Topic:     topic,
		Partition: 0, // 简单起见，先固定使用分区0
		Payload:   payload,
		Timestamp: time.Now(),
	}

	t, exists := p.broker.GetTopic(topic)
	if !exists {
		if err := p.broker.CreateTopic(topic, 1); err != nil {
			return err
		}
		t, _ = p.broker.GetTopic(topic)
	}

	partition := t.GetPartitions()[msg.Partition%len(t.GetPartitions())]
	msg.Offset = partition.AppendMessage(*msg)

	return nil
}
