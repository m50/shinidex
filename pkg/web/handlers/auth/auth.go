package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/m50/shinidex/pkg/database"
	"github.com/m50/shinidex/pkg/database/passwords"
	"github.com/m50/shinidex/pkg/views"
	"github.com/m50/shinidex/pkg/web/session"
	smiddleware "github.com/m50/shinidex/pkg/web/middleware"
)

func Router(e *echo.Echo) {
	group := e.Group("/auth")

	group.GET("/register", registerForm)
	group.POST("/register", register)

	group.GET("/login", loginForm)
	group.POST("/login", login, middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	group.POST("/logout", logout, smiddleware.AuthnMiddleware)
}

func registerForm(c echo.Context) error {
	return views.RenderView(c, http.StatusOK, RegisterForm())
}

func register(c echo.Context) error {
	return nil
}

func loginForm(c echo.Context) error {
	return views.RenderView(c, http.StatusOK, LoginForm())
}

func login(c echo.Context) error {
	ctx := c.(database.DBContext)
	email := c.FormValue("email")
	user, err := ctx.DB().Users().FindByEmail(email)
	if err != nil {
		return views.RenderError(c, err)
	}

	if err := passwords.ComparePasswords(user.Password, c.FormValue("password")); err != nil {
		return views.RenderViews(c, http.StatusForbidden, LoginForm(), views.Error(err))
	}

	if err := session.New(c, user); err != nil {
		c.Logger().Error(err)
		return views.RenderError(c, err)
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func logout(c echo.Context) error {
	session.Close(c)
	return c.Redirect(http.StatusMovedPermanently, "/auth/login")
}
