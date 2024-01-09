package messagqqueue

import (
	"github.com/adjust/rmq/v5"
	redisclient "github.com/dathuynh1108/clean-arch-base/pkg/redis_client"
	"github.com/dathuynh1108/clean-arch-base/pkg/singleton"
)

type MessageQueue interface {
	OpenQueue(queueName string) (rmq.Queue, error)
	GetQueue(queueName string) (rmq.Queue, error)
}

var (
	messageQueueSigleton *singleton.Singleton[MessageQueue]
)

func InitMessageQueue() error {
	messageQueueSigleton = singleton.NewSingleton[MessageQueue](func() MessageQueue {
		redisClient := redisclient.GetRedis()
		msgQueue, err := NewRedisMessageQueue(redisClient)
		if err != nil {
			panic(err)
		}
		return msgQueue
	})
	return nil
}

func GetMessageQueue() MessageQueue {
	return messageQueueSigleton.Get()
}
