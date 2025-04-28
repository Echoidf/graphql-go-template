package cmd

import (
	"gqlexample/graph"
	"net"
	"net/http"
	"os"

	"gqlexample/pkg/config"
	"gqlexample/pkg/middware"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
	"go.uber.org/zap"
)

var cfg = config.GetConfig()

func Run() {
	schema := graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver()})
	srv := handler.New(schema)

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AroundOperations(middware.GqlLogger)

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	socketPath := cfg.SocketPath

	// 确保socket文件不存在
	if err := os.RemoveAll(socketPath); err != nil {
		zap.L().Fatal("Failed to remove existing socket file", zap.Error(err))
	}

	// 创建 Unix Domain Socket 监听器
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		zap.L().Fatal("Failed to create unix socket", zap.Error(err))
	}
	defer listener.Close()

	// 设置socket文件权限
	if err := os.Chmod(socketPath, 0666); err != nil {
		zap.L().Fatal("Failed to chmod socket", zap.Error(err))
	}

	zap.L().Info("GraphQL server is running on unix socket",
		zap.String("socket", socketPath))

	if err := http.Serve(listener, nil); err != nil {
		zap.L().Fatal("Failed to start server", zap.Error(err))
	}
}
