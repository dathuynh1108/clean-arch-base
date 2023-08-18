package messagqqueue

import (
	"fmt"
	"os"
	"sync"

	"github.com/adjust/rmq/v5"
	"github.com/redis/go-redis/v9"
)

type RedisMessageQueue struct {
	mux        sync.RWMutex
	connection rmq.Connection
	queue      map[string]rmq.Queue
	errChan    chan error
}

func NewRedisMessageQueue(redisClient redis.UniversalClient) (*RedisMessageQueue, error) {
	errChan := make(chan<- error)
	connection, err := rmq.OpenConnectionWithRedisClient(os.Getenv("HOSTNAME"), redisClient, errChan)
	if err != nil {
		return nil, err
	}
	return &RedisMessageQueue{
		connection: connection,
		queue:      make(map[string]rmq.Queue),
	}, nil
}

func (q *RedisMessageQueue) OpenQueue(name string) (rmq.Queue, error) {
	q.mux.Lock()
	defer q.mux.Unlock()
	queue, err := q.connection.OpenQueue(name)
	if err != nil {
		return nil, err
	}
	q.queue[name] = queue
	return queue, nil
}

func (q *RedisMessageQueue) GetQueue(name string) (rmq.Queue, error) {
	q.mux.RLock()
	defer q.mux.RUnlock()
	queue, ok := q.queue[name]
	if !ok {
		return nil, fmt.Errorf("queue not found")
	}
	return queue, nil
}
