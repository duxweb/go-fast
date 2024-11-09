package queue

import (
	"context"
	"time"
)

type Queue interface {
	Worker(queueName string)
	Start() error
	Register(queueName, name string, callback func(ctx context.Context, params []byte) error) error
	Add(queueName string, add QueueAdd) (id string, err error)
	AddDelay(queueName string, add QueueAddDelay) (id string, err error)
	Del(queueName string, id string) error
	List(queueName string, page int, limit int) (data []QueueItem, count int64, err error)
	Names() []string
	Close() error
}

type QueueItem struct {
	ID        string         `json:"id"`
	QueueName string         `json:"queue_name"`
	Name      string         `json:"name"`
	Params    map[string]any `json:"params"`
	CreatedAt time.Time      `json:"created_at"`
	RunAt     time.Time      `json:"run_at"`
	Retried   int            `json:"retried"`
}

type QueueAdd struct {
	Name   string
	Params []byte
}

type QueueAddDelay struct {
	QueueAdd
	Delay time.Duration
}
