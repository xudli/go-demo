package main

import (
	"fmt"
	"kafka-demo/internal/broker"
	"kafka-demo/internal/consumer"
	"kafka-demo/internal/producer"
	"kafka-demo/pkg/types"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 创建broker实例
	b := broker.NewBroker()

	// 创建producer
	p := producer.NewProducer(b, types.Config{
		Brokers: []string{"localhost:9092"},
	})

	// 创建consumer
	c := consumer.NewConsumer(b, types.Config{
		Brokers:    []string{"localhost:9092"},
		GroupID:    "test-group",
		AutoCommit: true,
	})

	// 订阅主题
	messages, err := c.Subscribe("test-topic")
	if err != nil {
		fmt.Printf("Failed to subscribe: %v\n", err)
		return
	}

	// 设置信号处理
	done := make(chan bool)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// 在goroutine中接收消息
	go func() {
		for msg := range messages {
			fmt.Printf("Received message: %s\n", string(msg.Payload))
		}
		done <- true
	}()

	// 等待一秒确保消费者准备就绪
	time.Sleep(time.Second)
	fmt.Println("Consumer is ready")

	// 发送一些测试消息
	for i := 0; i < 5; i++ {
		message := fmt.Sprintf("Hello, World! Message %d", i)
		err := p.Send("test-topic", []byte(message))
		if err != nil {
			fmt.Printf("Failed to send message: %v\n", err)
			return
		}
		fmt.Printf("Sent message: %s\n", message)
		time.Sleep(time.Millisecond * 100)
	}

	// 等待消息处理完成或接收到终止信号
	select {
	case <-signals:
		fmt.Println("Received termination signal")
	case <-time.After(5 * time.Second):
		fmt.Println("Timeout waiting for messages")
	case <-done:
		fmt.Println("All messages processed")
	}
}
