package broker

import (
	"fmt"
	"kafka-demo/pkg/types"
	"sync"
)

type Broker struct {
	topics map[string]*Topic
	mu     sync.RWMutex
}

type Topic struct {
	Name       string
	Partitions []*Partition
	mu         sync.RWMutex
}

type Partition struct {
	ID  int
	Log *Log
	mu  sync.RWMutex
}

type Log struct {
	messages []types.Message
	mu       sync.RWMutex
}

func NewBroker() *Broker {
	return &Broker{
		topics: make(map[string]*Topic),
	}
}

// GetTopic 获取指定的Topic
func (b *Broker) GetTopic(name string) (*Topic, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	topic, exists := b.topics[name]
	return topic, exists
}

// GetMessages 获取指定分区的消息
func (p *Partition) GetMessages(offset int64) []types.Message {
	p.mu.RLock()
	defer p.mu.RUnlock()

	p.Log.mu.RLock()
	defer p.Log.mu.RUnlock()

	if offset >= int64(len(p.Log.messages)) {
		return nil
	}
	return p.Log.messages[offset:]
}

// GetPartitions 获取Topic的所有分区
func (t *Topic) GetPartitions() []*Partition {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.Partitions
}

// AppendMessage 添加消息到分区
func (p *Partition) AppendMessage(msg types.Message) int64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.Log.mu.Lock()
	defer p.Log.mu.Unlock()

	offset := int64(len(p.Log.messages))
	p.Log.messages = append(p.Log.messages, msg)
	fmt.Printf("Message appended to partition %d, offset: %d, total messages: %d\n",
		p.ID, offset, len(p.Log.messages))

	return offset
}

func (b *Broker) CreateTopic(name string, numPartitions int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, exists := b.topics[name]; exists {
		return nil
	}

	partitions := make([]*Partition, numPartitions)
	for i := 0; i < numPartitions; i++ {
		partitions[i] = &Partition{
			ID:  i,
			Log: &Log{messages: make([]types.Message, 0)},
		}
	}

	b.topics[name] = &Topic{
		Name:       name,
		Partitions: partitions,
	}

	return nil
}

func (b *Broker) Publish(msg *types.Message) error {
	b.mu.RLock()
	topic, exists := b.topics[msg.Topic]
	b.mu.RUnlock()

	if !exists {
		return b.CreateTopic(msg.Topic, 1)
	}

	partition := topic.Partitions[msg.Partition%len(topic.Partitions)]

	partition.mu.Lock()
	defer partition.mu.Unlock()

	partition.Log.mu.Lock()
	msg.Offset = int64(len(partition.Log.messages))
	partition.Log.messages = append(partition.Log.messages, *msg)
	partition.Log.mu.Unlock()

	return nil
}
