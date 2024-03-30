package views

import (
	"bytes"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/m50/shinidex/pkg/web/session"
)

func HeaderMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			return err
		}
		if c.Request().Header.Get("hx-request") != "true" {
			return nil
		}
		user, err := session.GetAuthedUser(c)
		if err != nil {
			c.Logger().Warn(err)
		}
		return AddView(c, Header(user))
	}
}

func renderWrappedView(ctx echo.Context, t templ.Component) error {
	user, err := session.GetAuthedUser(ctx)
	if err != nil {
		ctx.Logger().Warn(err)
	}
	base := BaseLayout(user)
	children := templ.WithChildren(ctx.Request().Context(), t)
	return base.Render(children, ctx.Response().Writer)
}

func RenderError(ctx echo.Context, err error) error {
	ctx.Logger().Error(err)
	return RenderView(ctx, http.StatusInternalServerError, Error(err))
}

func AddView(ctx echo.Context, t templ.Component) error {
	return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}

func RenderView(ctx echo.Context, status int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(status)
	if ctx.Request().Header.Get("hx-request") != "true" {
		return renderWrappedView(ctx, t)
	}
	return AddView(ctx, t)
}

func RenderViews(ctx echo.Context, status int, cmpts ...templ.Component) error {
	ctx.Response().Writer.WriteHeader(status)
	var err error
	for idx, cmpt := range cmpts {
		if idx == 0 && ctx.Request().Header.Get("hx-request") != "true" {
			err = renderWrappedView(ctx, cmpt)
		} else {
			err = AddView(ctx, cmpt)
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
