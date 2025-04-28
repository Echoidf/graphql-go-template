package middware

import (
	"context"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"go.uber.org/zap"
)

func GqlLogger(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	op := graphql.GetOperationContext(ctx)
	start := time.Now()

	resp := next(ctx)

	duration := time.Since(start)
	zap.L().Info("Gql operation",
		zap.String("operation", op.OperationName),
		// zap.String("query", op.RawQuery),
		zap.Any("variables", op.Variables),
		zap.Duration("duration", duration),
	)

	return resp
}
