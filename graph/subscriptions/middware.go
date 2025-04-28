package subscriptions

import "log"

// Middleware 订阅中间件接口
type Middleware interface {
	BeforeSubscribe(*Subscription) error
	AfterUnsubscribe(string)
}

// AuthMiddleware 认证中间件示例
type AuthMiddleware struct{}

func (m *AuthMiddleware) BeforeSubscribe(sub *Subscription) error {
	// 实现认证逻辑
	// if !isAuthorized(sub.Context, sub.Topic, sub.Channel) {
	//     return errors.New("unauthorized")
	// }
	return nil
}

func (m *AuthMiddleware) AfterUnsubscribe(id string) {
	// 清理资源
}

// LoggingMiddleware 日志中间件
type LoggingMiddleware struct{}

func (m *LoggingMiddleware) BeforeSubscribe(sub *Subscription) error {
	log.Printf("Subscription attempt: %v", sub)
	return nil
}

func (m *LoggingMiddleware) AfterUnsubscribe(id string) {
	log.Printf("Unsubscribed: %s", id)
}
