package views

import (
	"bytes"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gookit/slog"
	"github.com/labstack/echo/v4"
	"github.com/m50/shinidex/pkg/context"
	"github.com/m50/shinidex/pkg/web/session"
)

func renderWrappedView(ctx echo.Context, t templ.Component) error {
	user, _ := session.GetAuthedUser(ctx)
	v, ok := ctx.Get("rendersPokemon").(bool)
	rendersPkmn := ok && v
	base := BaseLayout(user, rendersPkmn)
	children := templ.WithChildren(context.FromEcho(ctx), t)
	return base.Render(children, ctx.Response().Writer)
}

func RenderErrorWithCode(ctx echo.Context, status int, err error) error {
	return RenderView(ctx, status, Error(err))
}

func RenderError(ctx echo.Context, err error) error {
	slog.WithContext(context.FromEcho(ctx)).Error(err)
	return RenderErrorWithCode(ctx, http.StatusInternalServerError, err)
}

func AddView(ctx echo.Context, t templ.Component) error {
	return t.Render(context.FromEcho(ctx), ctx.Response().Writer)
}

func RenderView(ctx echo.Context, status int, cmpts ...templ.Component) error {
	ctx.Response().WriteHeader(status)
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
	if err := t.Render(context.FromEcho(ctx), b); err != nil {
		return "", err
	}
	return b.String(), nil
}
