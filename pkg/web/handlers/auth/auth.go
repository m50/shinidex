package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/m50/shinidex/pkg/database"
	"github.com/m50/shinidex/pkg/database/passwords"
	"github.com/m50/shinidex/pkg/views"
)

type Context interface {
	DB() *database.Database
}

func Router(e *echo.Echo) {
	group := e.Group("/auth")
	group.GET("/", func (c echo.Context) error {
		return views.RenderView(c, http.StatusOK, LoginForm())
	})
	group.GET("/login", login, middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
}

func login(c echo.Context) error {
	ctx := c.(Context)
	email := c.FormValue("email")
	user, err := ctx.DB().Users().FindByEmail(email)
	if err != nil {
		return views.RenderError(c, err)
	}

	if err := passwords.ComparePasswords(user.Password, c.FormValue("password")); err != nil {
		return views.RenderViews(c, http.StatusForbidden, LoginForm(), views.Error(err))
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}