package messagqqueue

import (
	"github.com/adjust/rmq/v5"
	"github.com/dathuynh1108/clean-arch-base/pkg/redisclient"
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
	}, true)
	return nil
}

func GetMessageQueue() MessageQueue {
	return messageQueueSigleton.Get()
}
