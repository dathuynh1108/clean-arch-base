package jsonformater

import "context"

type ContextKey string

const (
	TraceContextKey ContextKey = "trace_key"
)

func WithTraceKey(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, TraceContextKey, key)
}

func getTraceKey(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if key, ok := ctx.Value(TraceContextKey).(string); ok {
		return key
	}

	return ""
}
