package consumer

import (
	"fmt"
	"kafka-demo/internal/broker"
	"kafka-demo/pkg/types"
	"sync"
	"time"
)

type Consumer struct {
	broker  *broker.Broker
	config  types.Config
	offsets map[string]map[int]int64
	mu      sync.RWMutex
}

func NewConsumer(broker *broker.Broker, config types.Config) *Consumer {
	return &Consumer{
		broker:  broker,
		config:  config,
		offsets: make(map[string]map[int]int64),
	}
}

func (c *Consumer) Subscribe(topic string) (<-chan types.Message, error) {
	messages := make(chan types.Message, 100)

	c.mu.Lock()
	if _, exists := c.offsets[topic]; !exists {
		c.offsets[topic] = make(map[int]int64)
	}
	c.mu.Unlock()

	go c.consumeMessages(topic, messages)

	return messages, nil
}

func (c *Consumer) consumeMessages(topic string, messages chan<- types.Message) {
	defer close(messages)
	fmt.Printf("Started consuming topic: %s\n", topic)

	for {
		t, exists := c.broker.GetTopic(topic)
		if !exists {
			fmt.Printf("Waiting for topic %s to be created...\n", topic)
			time.Sleep(time.Millisecond * 100)
			continue
		}

		c.processPartitions(topic, t.GetPartitions(), messages)
		time.Sleep(time.Millisecond * 100)
	}
}

func (c *Consumer) processPartitions(topic string, partitions []*broker.Partition, messages chan<- types.Message) {
	for _, partition := range partitions {
		c.mu.RLock()
		currentOffset := c.offsets[topic][partition.ID]
		c.mu.RUnlock()

		msgs := partition.GetMessages(currentOffset)
		if len(msgs) > 0 {
			fmt.Printf("Found %d messages in partition %d starting from offset %d\n",
				len(msgs), partition.ID, currentOffset)
		}

		for _, msg := range msgs {
			select {
			case messages <- msg:
				c.updateOffset(topic, partition.ID, currentOffset+1)
				fmt.Printf("Sent message to channel: %s\n", string(msg.Payload))
			default:
				fmt.Printf("Channel is full, message dropped\n")
			}
		}
	}
}

func (c *Consumer) updateOffset(topic string, partitionID int, offset int64) {
	c.mu.Lock()
	c.offsets[topic][partitionID] = offset
	c.mu.Unlock()
}
