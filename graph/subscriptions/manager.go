package subscriptions

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Manager struct {
	mu            sync.RWMutex
	subscriptions map[string]*Subscription
	eventChan     chan Event
	timeout       time.Duration
	middlewares   []Middleware
}

func NewManager(timeout time.Duration) *Manager {
	m := &Manager{
		subscriptions: make(map[string]*Subscription),
		eventChan:     make(chan Event, 100),
		timeout:       timeout,
	}

	go m.eventDispatcher()
	return m
}

// AddMiddleware 添加订阅中间件
func (m *Manager) AddMiddleware(mw Middleware) {
	m.middlewares = append(m.middlewares, mw)
}

// Subscribe 创建新订阅
func (m *Manager) Subscribe(ctx context.Context, topic SubscriptionTopic, channel string) (*Subscription, error) {
	sub := &Subscription{
		ID:      generateID(),
		Topic:   topic,
		Channel: channel,
		Output:  make(chan interface{}, 1),
		Context: ctx,
	}

	// 执行中间件链
	for _, mw := range m.middlewares {
		if err := mw.BeforeSubscribe(sub); err != nil {
			return nil, err
		}
	}

	m.mu.Lock()
	m.subscriptions[sub.ID] = sub
	m.mu.Unlock()

	// 启动清理协程
	go m.monitorSubscription(sub)

	log.Printf("New subscription: %s (%s/%s)", sub.ID, topic, channel)
	return sub, nil
}

// Unsubscribe 移除订阅
func (m *Manager) Unsubscribe(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if sub, exists := m.subscriptions[id]; exists {
		close(sub.Output)
		delete(m.subscriptions, id)
		log.Printf("Subscription removed: %s", id)
	}
}

// Publish 发布事件
func (m *Manager) Publish(event Event) {
	select {
	case m.eventChan <- event:
	default:
		log.Println("Event channel full, dropping event")
	}
}

// 事件分发器
func (m *Manager) eventDispatcher() {
	for event := range m.eventChan {
		m.mu.RLock()

		for _, sub := range m.subscriptions {
			if sub.Topic == event.Topic && sub.Channel == event.Channel {
				select {
				case sub.Output <- event.Payload:
				case <-time.After(m.timeout):
					log.Printf("Timeout sending to subscriber %s", sub.ID)
				}
			}
		}

		m.mu.RUnlock()
	}
}

// 监控订阅状态
func (m *Manager) monitorSubscription(sub *Subscription) {
	<-sub.Context.Done()
	m.Unsubscribe(sub.ID)
}

func generateID() string {
	// 实现一个唯一的ID生成逻辑
	return "sub_" + time.Now().Format("20060102150405") + "_" + randString(6)
}

func randString(n int) string {
	// 简化实现，实际项目使用更安全的随机生成
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
