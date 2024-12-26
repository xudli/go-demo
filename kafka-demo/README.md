# Go简易消息队列实现

## 项目介绍

这是一个用Go语言实现的轻量级消息队列系统，实现了基本的发布/订阅功能。该项目主要用于学习消息队列的基本原理和实现方式。

## 核心特性

- 基于内存的消息存储
- 发布/订阅模式
- 多主题（Topic）支持
- 分区（Partition）机制
- 基于 offset 的消息消费
- 并发安全的消息处理

## 系统架构

### 核心组件

1. **Broker**
   - 负责消息的存储和管理
   - 管理 Topic 和 Partition
   - 提供消息读写接口
   - 实现并发控制

2. **Producer**
   - 发送消息到指定 Topic
   - 支持自动创建 Topic
   - 消息自动分配到 Partition

3. **Consumer**
   - 订阅指定 Topic
   - 维护消费位置（offset）
   - 支持等待 Topic 创建
   - 异步消息消费

## 数据结构

### Message

```go
type Message struct {
    ID        string
    Topic     string
    Partition int
    Payload   []byte
    Timestamp time.Time
    Offset    int64
}
```

### Broker

```go
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
```

## 实现细节

### 消息发布流程

1. Producer 调用 Send 方法发送消息
2. 如果 Topic 不存在，自动创建
3. 消息被追加到指定分区的日志中
4. 返回消息的 offset

### 消息订阅流程

1. Consumer 调用 Subscribe 方法订阅 Topic
2. 如果 Topic 不存在，等待 Topic 创建
3. 从上次消费的 offset 开始读取消息
4. 通过 channel 异步发送消息给消费者

### 并发控制

- 使用读写锁（sync.RWMutex）保护共享资源
- Broker 级别的锁保护 topics 映射
- Topic 级别的锁保护 Partitions
- Partition 级别的锁保护消息日志
- Consumer 级别的锁保护 offset 映射

## 使用示例

### 发送消息

```go
producer := kafka.NewProducer(broker, types.Config{
    Brokers: []string{"localhost:9092"},
})

err := producer.Send("test-topic", []byte("Hello World"))
```

### 消费消息

```go
consumer := kafka.NewConsumer(broker, types.Config{
    Brokers:    []string{"localhost:9092"},
    GroupID:    "test-group",
    AutoCommit: true,
})

messages, err := consumer.Subscribe("test-topic")
for msg := range messages {
    fmt.Printf("Received: %s\n", string(msg.Payload))
}
```

## 限制说明

1. 仅支持内存存储，不支持持久化
2. 不支持消息确认机制
3. 不支持消费者组
4. 简单的分区分配策略（固定使用分区0）
5. 不支持消息压缩

## 后续改进方向

1. 添加消息持久化
2. 实现消费者组功能
3. 添加消息确认机制
4. 优化分区分配策略
5. 添加消息压缩支持
6. 实现更高效的消息传递机制
