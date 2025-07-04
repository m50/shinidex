package middleware

import (
	"github.com/gookit/slog"
	"github.com/labstack/echo/v4"
	"github.com/m50/shinidex/pkg/context"
	"github.com/m50/shinidex/pkg/views"
	"github.com/m50/shinidex/pkg/web/errors"
)

func ErrorHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		e := next(c)
		switch err := e.(type) {
		case errors.APIError:
			slog.WithContext(context.FromEcho(c)).Error(err)
			return views.RenderErrorWithCode(c, err.StatusCode, err)
		case nil:
			return nil
		default:
			slog.WithContext(context.FromEcho(c)).Error(err)
			return err
		}
	}
}
