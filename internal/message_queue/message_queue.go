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
	MessageQueueSigleton *singleton.Singleton[MessageQueue]
)

func InitMessageQueue() {
	MessageQueueSigleton = singleton.NewSingleton[MessageQueue](func() MessageQueue {
		redisClient := redisclient.GetRedis()
		msgQueue, err := NewRedisMessageQueue(redisClient)
		if err != nil {
			panic(err)
		}
		return msgQueue
	})
}

func GetMessageQueue() MessageQueue {
	return MessageQueueSigleton.Get()
}
