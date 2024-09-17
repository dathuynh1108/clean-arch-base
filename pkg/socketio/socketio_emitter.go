package socketio

import (
	"sync"

	"github.com/dathuynh1108/clean-arch-base/pkg/redisclient"

	emitter "github.com/dathuynh1108/go-socketio-redis-emitter"
)

func GetEmitter() *emitter.Emitter {
	return emitterPool.Get().(*emitter.Emitter)
}

func PutEmitter(e *emitter.Emitter) {
	e.Reset()
	emitterPool.Put(e)
}

var emitterPool = sync.Pool{
	New: func() any {
		return emitter.NewEmitter(&emitter.Options{
			Key:   "SOCKETIO KEY",
			Redis: redisclient.GetRedis(),
		})
	},
}
