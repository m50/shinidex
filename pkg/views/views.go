package views

import (
	"bytes"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func renderWrappedView(ctx echo.Context, t templ.Component) error {
	base := BaseLayout()
	children := templ.WithChildren(ctx.Request().Context(), t)
	return base.Render(children, ctx.Response().Writer)
}

func RenderError(ctx echo.Context, err error) error {
	return RenderView(ctx, http.StatusInternalServerError, Error(err))
}

func RenderView(ctx echo.Context, status int, t templ.Component) error {
    ctx.Response().Writer.WriteHeader(status)
	if ctx.Request().Header.Get("hx-request") != "true" {
		return renderWrappedView(ctx, t)
	}
    return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}

func RenderViews(ctx echo.Context, status int, cmpts ...templ.Component) error {
    ctx.Response().Writer.WriteHeader(status)
	var err error
	for idx, cmpt := range cmpts {
		if idx == 0 && ctx.Request().Header.Get("hx-request") != "true" {
			err = renderWrappedView(ctx, cmpt)
		} else {
			err = cmpt.Render(ctx.Request().Context(), ctx.Response().Writer)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func ToStr(ctx echo.Context, t templ.Component) (string, error) {
	b := new(bytes.Buffer)
	if err := t.Render(ctx.Request().Context(), b); err != nil {
		return "", err
	}
	return b.String(), nil
}
