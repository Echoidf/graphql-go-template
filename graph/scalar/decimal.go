package scalar

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
	"github.com/shopspring/decimal"
)

// MarshalDecimal 将 decimal.Decimal 转换为 GraphQL 标量
func MarshalDecimal(d decimal.Decimal) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, d.String())
	})
}

// UnmarshalDecimal 将输入值转换为 decimal.Decimal
func UnmarshalDecimal(v interface{}) (decimal.Decimal, error) {
	switch v := v.(type) {
	case string:
		return decimal.NewFromString(v)
	case float64:
		return decimal.NewFromFloat(v), nil
	case int:
		return decimal.NewFromInt(int64(v)), nil
	case json.Number:
		return decimal.NewFromString(v.String())
	default:
		return decimal.Zero, fmt.Errorf("无法将 %T 转换为 Decimal", v)
	}
}
