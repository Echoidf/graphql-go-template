package event

import (
	"gqlexample/pkg/cache"
	"maps"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

type DataEvent struct {
	Data  any
	Topic string
}

type DataChannel chan DataEvent

// EventBus 存储有关订阅者感兴趣的特定主题的信息
type EventBus struct {
	subscribers *cache.Cache[string, map[int32]DataChannel]
	msgQueues   *cache.Cache[string, DataChannel]
	idCounter   int32
	mu          sync.RWMutex
}

type SubParam struct {
	Topic      string
	Ch         DataChannel
	BufferSize uint
	Timeout    *time.Duration
}

type SubOption func(*SubParam)

func WithBufferSize(size uint) SubOption {
	return func(p *SubParam) {
		p.BufferSize = size
	}
}

func WithTimeout(timeout time.Duration) SubOption {
	return func(p *SubParam) {
		p.Timeout = &timeout
	}
}

const DEFAULT_BUFFER_SIZE = 1024

var Eb *EventBus

func init() {
	Eb = NewEventBus()
}

// NewEventBus 创建一个新的事件总线
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: cache.NewCache[string, map[int32]DataChannel](),
		msgQueues:   cache.NewCache[string, DataChannel](),
	}
}

// Subscribe 订阅特定主题
func (eb *EventBus) Subscribe(topic string, ch DataChannel, opts ...SubOption) int32 {
	atomic.AddInt32(&eb.idCounter, 1)
	subId := eb.idCounter

	subParam := &SubParam{
		Topic:      topic,
		Ch:         ch,
		BufferSize: DEFAULT_BUFFER_SIZE,
	}

	for _, opt := range opts {
		opt(subParam)
	}

	eb.subscribers.Update(topic, func(m *map[int32]DataChannel) error {
		if m == nil || *m == nil {
			*m = make(map[int32]DataChannel)
			eb.msgQueues.Set(topic, make(DataChannel, subParam.BufferSize))
			go eb.processQueue(subParam)
		}
		(*m)[subId] = ch
		return nil
	})

	return subId
}

// processQueue 处理特定主题的消息队列
func (eb *EventBus) processQueue(subParam *SubParam) {
	queue, _ := eb.msgQueues.Get(subParam.Topic)
	var subscribers map[int32]DataChannel
	var timeoutChan <-chan time.Time
	if subParam.Timeout != nil {
		timeoutChan = time.After(*subParam.Timeout)
	}

	defer func() {
		if r := recover(); r != nil {
			zap.L().Warn("processQueue panic", zap.Any("panic", r))
		}
	}()

	for {
		select {
		case event, ok := <-queue:
			if !ok {
				return
			}
			// 获取最新的订阅者列表
			eb.mu.RLock()
			subers, _ := eb.subscribers.Get(subParam.Topic)
			subscribers = make(map[int32]DataChannel, len(subers))
			maps.Copy(subscribers, subers)
			eb.mu.RUnlock()

			// 向所有订阅者发送消息
			for _, ch := range subscribers {
				select {
				case ch <- event:
				case <-timeoutChan:
					zap.L().Warn("processQueue timeout", zap.Any("topic", subParam.Topic))
					return
				}
			}
		}
	}
}

// Publish 发布消息到特定主题
func (eb *EventBus) Publish(topic string, data any) {
	eb.mu.RLock()
	queue, exists := eb.msgQueues.Get(topic)
	eb.mu.RUnlock()

	defer func() {
		if r := recover(); r != nil {
			zap.L().Warn("processQueue panic", zap.Any("panic", r))
		}
	}()

	if exists {
		select {
		case queue <- DataEvent{data, topic}:
		}
	}
}

// Unsubscribe 取消订阅特定主题
func (eb *EventBus) Unsubscribe(topic string, subId int32) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if chans, found := eb.subscribers.Get(topic); found {
		close(chans[subId])
		delete(chans, subId)
		eb.subscribers.Set(topic, chans)

		if len(chans) == 0 {
			queue, _ := eb.msgQueues.Get(topic)
			close(queue)
			eb.msgQueues.Delete(topic)
			eb.subscribers.Delete(topic)
		}
	}
}
