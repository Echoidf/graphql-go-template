package subscriptions

import (
	"context"
)

// SubscriptionTopic 订阅主题标识
type SubscriptionTopic string

const (
	TopicMessages SubscriptionTopic = "messages"
	TopicUsers    SubscriptionTopic = "users"
)

// Subscription 表示一个活跃的订阅
type Subscription struct {
	ID      string
	Topic   SubscriptionTopic
	Channel string
	Output  chan any
	Context context.Context
}

// Event 发布的事件结构
type Event struct {
	Topic   SubscriptionTopic
	Channel string
	Payload any
}
