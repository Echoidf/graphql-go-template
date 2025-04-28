package graph

//go:generate go run github.com/99designs/gqlgen generate
import (
	"gqlexample/graph/model"
	"gqlexample/graph/subscriptions"
	"time"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	todos               []*model.Todo
	SubscriptionManager *subscriptions.Manager
}

func NewResolver() *Resolver {
	// 初始化订阅管理器，设置10秒超时
	mgr := subscriptions.NewManager(10 * time.Second)

	// 添加中间件
	// mgr.AddMiddleware(&subscriptions.AuthMiddleware{})
	// mgr.AddMiddleware(&subscriptions.LoggingMiddleware{})

	return &Resolver{
		SubscriptionManager: mgr,
	}
}
