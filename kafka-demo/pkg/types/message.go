package types

import "time"

type Message struct {
	ID        string
	Topic     string
	Partition int
	Payload   []byte
	Timestamp time.Time
	Offset    int64
}

type Config struct {
	Brokers    []string
	GroupID    string
	AutoCommit bool
}
