package context

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
)

type ck string

var ckEchoContext = ck("echo.Context")

type echoContext struct {
	context.Context
	echoCtx echo.Context
}

func FromEcho(ctx echo.Context) context.Context {
	return &echoContext{
		Context: ctx.Request().Context(),
		echoCtx: ctx,
	}
}

func (c echoContext) Value(key any) any {
	if key == ckEchoContext {
		return c.Context
	}
	if k, ok := key.(string); ok {
		v := c.echoCtx.Get(k)
		return v
	}
	return c.Context.Value(key)
}

func (c *echoContext) String() string {
	return fmt.Sprintf("%s.FromEcho(echo.Context)", c.Context)
}

func ToEcho(ctx context.Context) echo.Context {
	if ctx == nil {
		return nil
	}
	if c, ok := ctx.(*echoContext); ok {
		return c.echoCtx
	}
	if c, ok := ctx.Value("echo.Context").(echo.Context); ok {
		return c
	}
	return nil
}
