package common

import (
	"context"
	"time"

	"github.com/dathuynh1108/clean-arch-base/pkg/comjson"
	"github.com/labstack/echo/v4"
)

type EchoWrappedContext struct {
	echo.Context
	ctx context.Context
}

var _ context.Context = (*EchoWrappedContext)(nil)

func EchoWrapContext(c echo.Context) *EchoWrappedContext {
	if wrappedCtx, ok := c.(*EchoWrappedContext); ok {
		return wrappedCtx
	}
	return EchoWrapWithContext(c, c.Request().Context())
}

func EchoWrapWithContext(c echo.Context, ctx context.Context) *EchoWrappedContext {
	return &EchoWrappedContext{
		Context: c,
		ctx:     ctx,
	}
}

func (c *EchoWrappedContext) Deadline() (time.Time, bool) {
	return c.ctx.Deadline()
}

func (c *EchoWrappedContext) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *EchoWrappedContext) Err() error {
	return c.ctx.Err()
}

func (c *EchoWrappedContext) Value(key interface{}) interface{} {
	if keyStr, ok := key.(string); ok {
		if result := c.Get(keyStr); result != nil {
			return result
		}
	}
	return c.ctx.Value(key)
}

func (c *EchoWrappedContext) SetValue(key interface{}, value interface{}) {
	switch keyT := key.(type) {
	case string:
		c.Set(keyT, value)
	default:
		c.ctx = context.WithValue(c.ctx, key, value)
	}
}

func (c *EchoWrappedContext) StackUp(fn func(context.Context) context.Context) {
	c.ctx = fn(c.ctx)
}

// EchoJSONSerializer implements JSON encoding using encoding/json.
type EchoJSONSerializer struct{}

// Serialize converts an interface into a json and writes it to the response.
// You can optionally use the indent parameter to produce pretty JSONs.
func (d EchoJSONSerializer) Serialize(c echo.Context, i interface{}, indent string) error {
	enc := comjson.NewEncoder(c.Response())
	if indent != "" {
		enc.SetIndent("", indent)
	}
	return enc.Encode(i)
}

// Deserialize reads a JSON from a request body and converts it into an interface.
func (d EchoJSONSerializer) Deserialize(c echo.Context, i interface{}) error {
	return comjson.NewDecoder(c.Request().Body).Decode(i)
}
