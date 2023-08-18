package messagqqueue

import (
	"sync/atomic"

	"github.com/adjust/rmq/v5"
)

type Consumer struct {
	handler atomic.Value
}

type BatchConsumer struct {
	handler atomic.Value
}

func NewConsumer(f func(payload string) error) *Consumer {
	c := &Consumer{}
	c.handler.Store(f)
	return c
}

func NewBatchConsumer(f func(payload string) error) *BatchConsumer {
	c := &BatchConsumer{}
	c.handler.Store(f)
	return c
}

func (c *Consumer) Consume(delivery rmq.Delivery) {
	handler, ok := c.handler.Load().(func(payload string) error)
	if !ok {
		if err := delivery.Reject(); err != nil {
		}
		return
	}
	err := handler(delivery.Payload())
	if err != nil {
		if err := delivery.Reject(); err != nil {
		}
		return
	}
	if err := delivery.Ack(); err != nil {
		return
	}
}

func (c *BatchConsumer) Consume(deliveries rmq.Deliveries) {
	handler, ok := c.handler.Load().(func(payload []string) error)
	if !ok {
		if err := deliveries.Reject(); err != nil {
		}
		return
	}
	err := handler(deliveries.Payloads())
	if err != nil {
		if err := deliveries.Reject(); err != nil {
		}
		return
	}
	if err := deliveries.Ack(); err != nil {
		return
	}
}
