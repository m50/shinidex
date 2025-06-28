package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gookit/slog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/m50/shinidex/pkg/context"
	"github.com/m50/shinidex/pkg/database"
	"github.com/m50/shinidex/pkg/database/passwords"
	"github.com/m50/shinidex/pkg/types"
	"github.com/m50/shinidex/pkg/views"
	smiddleware "github.com/m50/shinidex/pkg/web/middleware"
	"github.com/m50/shinidex/pkg/web/session"
	"github.com/spf13/viper"
)

func Router(e *echo.Echo) {
	group := e.Group("/auth")

	if !viper.GetBool("auth.disable-registration") {
		group.GET("/register", registerForm)
		group.POST("/register", register, middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	}

	group.GET("/login", loginForm)
	group.POST("/login", login, middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	group.POST("/logout", logout, smiddleware.AuthnMiddleware)
}

type registerFormData struct {
	Email           string `form:"email"`
	Password        string `form:"password"`
	ConfirmPassword string `form:"confirm_password"`
	Honeypot        string `form:"name"`
}

func registerForm(c echo.Context) error {
	return views.RenderView(c, http.StatusOK, RegisterForm(registerFormData{}))
}

func register(c echo.Context) error {
	ctx := context.FromEcho(c)
	var f registerFormData
	if err := c.Bind(&f); err != nil {
		return err
	}
	if f.Honeypot != "" {
		return c.Redirect(http.StatusMovedPermanently, "/")
	}
	if f.Password != f.ConfirmPassword {
		err := errors.New("password and confirm password must be the same")
		return views.RenderView(c, http.StatusUnprocessableEntity,
			RegisterForm(f), views.Error(err))
	}

	db := c.(database.DBContext).DB()
	p, err := passwords.HashPassword(f.Password)
	if err != nil {
		return err
	}
	u := types.User{
		Email:    f.Email,
		Password: p,
	}
	userID, err := db.Users().Insert(ctx, u)
	if err != nil {
		return err
	}
	user, err := db.Users().FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := session.New(c, user); err != nil {
		slog.WithContext(ctx).Error(err)
		return views.RenderError(c, err)
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func loginForm(c echo.Context) error {
	return views.RenderView(c, http.StatusOK, LoginForm())
}

func login(c echo.Context) error {
	ctx := context.FromEcho(c)
	db := c.(database.DBContext).DB()
	email := c.FormValue("email")
	user, err := db.Users().FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return views.RenderView(c, http.StatusInternalServerError, LoginForm(),
			views.Error(err))
	} else if errors.Is(err, sql.ErrNoRows) {
		return views.RenderView(c, http.StatusForbidden, LoginForm(),
			views.Error(fmt.Errorf("no account found for %s", email)))
	}

	if err := passwords.ComparePasswords(user.Password, c.FormValue("password")); err != nil {
		slog.WithContext(ctx).Error(err)
		return views.RenderView(c, http.StatusForbidden, LoginForm(), views.Error(err))
	}

	if err := session.New(c, user); err != nil {
		slog.WithContext(ctx).Error(err)
		return views.RenderView(c, http.StatusInternalServerError, LoginForm(),
			views.Error(err))
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func logout(c echo.Context) error {
	session.Close(c)
	return c.Redirect(http.StatusMovedPermanently, "/auth/login")
}
